package inspector

import (
	"context"
	"embed"
	"html/template"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/types"
)

//go:embed html
var html embed.FS

var t = template.Must(template.ParseFS(html, "html/*"))

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type MessageDirection string

const (
	MessageDirectionRequest  MessageDirection = "request"
	MessageDirectionResponse MessageDirection = "response"
)

// MessageInfo represents a single MCP message for inspection
type MessageInfo struct {
	Timestamp string           `json:"timestamp"`
	Direction MessageDirection `json:"direction"`
	Content   string           `json:"content"`
}

type Inspector struct {
	listenAddress string
	messageChan   chan MessageInfo
	clients       map[*websocket.Conn]bool
	mutex         sync.RWMutex
	logger        types.Logger
	server        *http.Server
	isClosing     bool
}

func NewInspector(config *config.InspectorInfo, logger types.Logger) *Inspector {
	return &Inspector{
		listenAddress: config.ListenAddress,
		messageChan:   make(chan MessageInfo, 100), // Buffer of 100 messages
		clients:       make(map[*websocket.Conn]bool),
		logger:        logger,
		server:        nil,
		isClosing:     false,
	}
}

// EnqueueMessage adds a new message to the inspection queue
func (i *Inspector) EnqueueMessage(msg MessageInfo) {
	// json encode the message
	select {
	case i.messageChan <- msg:
		// Message enqueued successfully
	default:
		// Channel is full, message is dropped
	}
}

func (i *Inspector) Start(ctx context.Context) error {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})
	router.HandleFunc("/ws", i.serveWs)

	server := &http.Server{
		Addr:    i.listenAddress,
		Handler: router,
	}
	i.server = server

	errChan := make(chan error, 1)

	// Start the server in a goroutine
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if i.isClosing {
				return
			}
			i.logger.Error("error starting inspector", types.LogArg{
				"error": err,
			})
			if !i.isClosing {
				errChan <- err
			}
		}
	}()

	// Start the message broadcaster in a goroutine
	go func() {
		i.broadcastMessages()
		errChan <- nil
	}()

	select {
	case err := <-errChan:
		i.Close(ctx)
		return err
	case <-ctx.Done():
		i.Close(ctx)
		return ctx.Err()
	}
}

// broadcastMessages continuously reads from messageChan
// and broadcasts to all connected websocket clients
func (i *Inspector) broadcastMessages() {
	for msg := range i.messageChan {
		i.mutex.RLock()
		for client := range i.clients {
			err := client.WriteJSON(msg)
			if err != nil {
				client.Close()
				delete(i.clients, client)
			}
		}
		i.mutex.RUnlock()
	}
}

func (i *Inspector) serveWs(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		i.logger.Error("upgrade:", types.LogArg{
			"error": err,
		})
		return
	}
	defer c.Close()

	i.mutex.Lock()
	i.clients[c] = true
	i.mutex.Unlock()

	defer func() {
		i.mutex.Lock()
		delete(i.clients, c)
		i.mutex.Unlock()
	}()

	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			i.logger.Error("read:", types.LogArg{
				"error": err,
			})
			break
		}
	}
}

func (i *Inspector) shutdown() {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	for client := range i.clients {
		client.Close()
	}
	close(i.messageChan)
}

func (i *Inspector) Close(ctx context.Context) {
	i.logger.Info("Shutting down inspector", types.LogArg{
		"listenAddress": i.listenAddress,
	})

	if i.isClosing {
		return
	}
	i.isClosing = true

	if i.server != nil {
		// shutdown the server
		if err := i.server.Shutdown(ctx); err != nil {
			i.logger.Error("error shutting down inspector", types.LogArg{
				"error": err,
			})
		}
	}

	// shutdown the inspector
	i.shutdown()

}
