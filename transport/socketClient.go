package transport

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/llmcontext/gomcp/types"
)

// SocketClientTransport implements the Transport interface using TCP sockets
type SocketClientTransport struct {
	addr string
	conn net.Conn

	// Callback functions
	onMessage func(json.RawMessage)
	onClose   func()
	onError   func(error)

	// Mutex for thread-safe operations
	mu sync.Mutex

	// Channel to signal shutdown
	done chan struct{}
}

// NewSocketTransport creates a new socket transport instance
func NewSocketClientTransport(address string) types.Transport {
	return &SocketClientTransport{
		addr: address,
		done: make(chan struct{}),
	}
}

// Start implements Transport.Start
func (s *SocketClientTransport) Start(ctx context.Context) error {
	var err error
	s.conn, err = net.Dial("tcp", s.addr)
	if err != nil {
		return err
	}

	// Start reading messages in a separate goroutine
	go s.readLoop()

	return nil
}

// Send implements Transport.Send
func (s *SocketClientTransport) Send(message json.RawMessage) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		return fmt.Errorf("connection is closed")
	}

	// Add newline as message delimiter
	message = append(message, '\n')
	_, err := s.conn.Write(message)
	return err
}

// OnMessage implements Transport.OnMessage
func (s *SocketClientTransport) OnMessage(callback func(json.RawMessage)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onMessage = callback
}

// OnClose implements Transport.OnClose
func (s *SocketClientTransport) OnClose(callback func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onClose = callback
}

// OnError implements Transport.OnError
func (s *SocketClientTransport) OnError(callback func(error)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.onError = callback
}

// Close implements Transport.Close
func (s *SocketClientTransport) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}

	// Signal shutdown to readLoop
	close(s.done)

	if s.onClose != nil {
		s.onClose()
	}
}

// readLoop continuously reads messages from the socket
func (s *SocketClientTransport) readLoop() {
	reader := bufio.NewReader(s.conn)

	for {
		select {
		case <-s.done:
			return
		default:
			// Read until newline
			line, err := reader.ReadString('\n')
			if err != nil {
				if s.onError != nil {
					s.onError(err)
				}
				s.Close()
				return
			}

			// Parse the JSON from the line
			var message json.RawMessage
			if err := json.Unmarshal([]byte(line), &message); err != nil {
				if s.onError != nil {
					s.onError(err)
				}
				continue
			}

			if s.onMessage != nil {
				s.onMessage(message)
			}
		}
	}
}