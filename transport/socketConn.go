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

type SocketConn struct {
	conn net.Conn

	// Mutex for thread-safe operations
	mu sync.Mutex

	// Callback functions
	onMessage func(json.RawMessage)
	onClose   func()
	onError   func(error)
}

func NewSocketConn(conn net.Conn) types.Transport {
	return &SocketConn{
		conn: conn,
		mu:   sync.Mutex{},
	}
}

func (s *SocketConn) Start(ctx context.Context) (chan error, error) {
	errChan := make(chan error, 1)

	go s.readLoop(ctx, errChan)
	return errChan, nil
}

// Send implements Transport.Send
func (s *SocketConn) Send(message json.RawMessage) error {
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
func (s *SocketConn) OnMessage(callback func(json.RawMessage)) {
	s.onMessage = callback
}

// OnClose implements Transport.OnClose
func (s *SocketConn) OnClose(callback func()) {
	s.onClose = callback
}

// OnError implements Transport.OnError
func (s *SocketConn) OnError(callback func(error)) {
	s.onError = callback
}

// Close implements Transport.Close
func (s *SocketConn) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}

	// report the close
	if s.onClose != nil {
		s.onClose()
	}
}

// readLoop continuously reads messages from the socket
func (s *SocketConn) readLoop(ctx context.Context, errChan chan error) {
	reader := bufio.NewReader(s.conn)

	for {
		select {
		case <-ctx.Done():
			// the context is done, we should exit
			// we close the connection
			s.Close()
			return
		default:
			// Read until newline
			line, err := reader.ReadString('\n')
			if err != nil {
				if s.onError != nil {
					s.onError(err)
				}
				s.Close()
				// send the error to the errChan
				errChan <- err
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
