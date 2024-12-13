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

type SocketServer struct {
	address     string
	listener    net.Listener
	conn        net.Conn
	mutex       sync.Mutex
	onMessage   func(json.RawMessage)
	onClose     func()
	onError     func(error)
	isConnected bool
}

func NewSocketServer(address string) types.Transport {
	return &SocketServer{
		address:  address,
		listener: nil,
	}
}

func (s *SocketServer) Start(ctx context.Context) error {
	go func() {
		listener, err := net.Listen("tcp", s.address)
		if err != nil {
			if s.onError != nil {
				s.onError(err)
			}
			return
		}
		s.listener = listener
		for {
			select {
			case <-ctx.Done():
				s.Close()
				return
			default:
				conn, err := s.listener.Accept()
				if err != nil {
					if s.onError != nil {
						s.onError(err)
					}
					continue
				}

				s.mutex.Lock()
				s.conn = conn
				s.isConnected = true
				s.mutex.Unlock()

				go s.handleConnection(conn)
			}
		}
	}()

	return nil
}

func (s *SocketServer) handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		s.mutex.Lock()
		s.isConnected = false
		s.mutex.Unlock()
		if s.onClose != nil {
			s.onClose()
		}
	}()

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := scanner.Bytes()
		if s.onMessage != nil {
			s.onMessage(json.RawMessage(line))
		}
	}

	if err := scanner.Err(); err != nil {
		if s.onError != nil {
			s.onError(err)
		}
	}
}

func (s *SocketServer) Send(message json.RawMessage) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if !s.isConnected {
		return fmt.Errorf("not connected")
	}

	return json.NewEncoder(s.conn).Encode(message)
}

func (s *SocketServer) OnMessage(callback func(json.RawMessage)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.onMessage = callback
}

func (s *SocketServer) OnClose(callback func()) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.onClose = callback
}

func (s *SocketServer) OnError(callback func(error)) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.onError = callback
}

func (s *SocketServer) Close() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.conn != nil {
		s.conn.Close()
	}
	if s.listener != nil {
		s.listener.Close()
	}
	s.isConnected = false
}
