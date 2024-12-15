package socket

import (
	"net"

	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type SocketServer struct {
	address  string
	listener net.Listener
	onError  func(error)
}

func NewSocketServer(address string) *SocketServer {
	return &SocketServer{
		address:  address,
		listener: nil,
	}
}

func (s *SocketServer) OnError(callback func(error)) {
	s.onError = callback
}

func (s *SocketServer) Start(callback func(types.Transport)) error {
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
			conn, err := s.listener.Accept()
			if err != nil {
				if s.onError != nil {
					s.onError(err)
				}
				continue
			}

			transport := transport.NewSocketConn(conn)
			callback(transport)
		}
	}()

	return nil
}

func (s *SocketServer) Close() {
	if s.listener != nil {
		s.listener.Close()
	}
}
