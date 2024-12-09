package inspector

import (
	"embed"
	"html/template"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/llmcontext/gomcp/config"
)

//go:embed html
var html embed.FS

var t = template.Must(template.ParseFS(html, "html/*"))

type MessageDirection string

const (
	MessageDirectionRequest  MessageDirection = "request"
	MessageDirectionResponse MessageDirection = "response"
)

// MessageInfo represents a single MCP message for inspection
type MessageInfo struct {
	Timestamp string
	Direction MessageDirection
	Content   string
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
	select {
	case i.messageChan <- msg:
		// Message enqueued successfully
	default:
		// Channel is full, message is dropped
	}
}

func (inspector *Inspector) StartInspector() *Inspector {
	router := http.NewServeMux()
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if err := t.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	})

	server := &http.Server{
		Addr:    inspector.listenAddress,
		Handler: router,
	}

	// Start the server in a goroutine
	go func() {
		server.ListenAndServe()
	}()

	// Start the message broadcaster in a goroutine
	go inspector.broadcastMessages()

	return inspector
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
