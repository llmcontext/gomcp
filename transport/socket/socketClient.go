package socket

import (
	"net"

	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

// SocketClient implements the Transport interface using TCP sockets
type SocketClient struct {
	addr string
	conn net.Conn
}

// NewSocketTransport creates a new socket transport instance
func NewSocketClient(address string) *SocketClient {
	return &SocketClient{
		addr: address,
	}
}

func (s *SocketClient) Start() (types.Transport, error) {
	var err error
	s.conn, err = net.Dial("tcp", s.addr)
	if err != nil {
		return nil, err
	}

	conn := transport.NewSocketConn(s.conn)

	return conn, nil
}
