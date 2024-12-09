package inspector

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/logger"
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
}

func NewInspector(config *config.InspectorInfo) *Inspector {
	return &Inspector{
		listenAddress: config.ListenAddress,
		messageChan:   make(chan MessageInfo, 100), // Buffer of 100 messages
		clients:       make(map[*websocket.Conn]bool),
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

func (inspector *Inspector) StartInspector(ctx context.Context) error {
	router := http.NewServeMux()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})
	router.HandleFunc("/ws", inspector.serveWs)

	server := &http.Server{
		Addr:    inspector.listenAddress,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		server.ListenAndServe()
	}()

	// Start the message broadcaster in a goroutine
	go func() {
		inspector.broadcastMessages()
	}()

	// wait for the context to be done
	<-ctx.Done()

	fmt.Printf("# [inspector] shutdown\n")

	// shutdown the server
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("error shutting down inspector: %s\n", err)
	}

	// shutdown the inspector
	inspector.shutdown()

	return nil
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
		logger.Error("upgrade:", logger.Arg{
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
			logger.Error("read:", logger.Arg{
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
