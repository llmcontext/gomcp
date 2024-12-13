package mux

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type Multiplexer struct {
	listenAddress string
	transport     types.Transport
}

func NewMultiplexer(config *config.ProxyConfig) *Multiplexer {
	return &Multiplexer{
		listenAddress: config.ListenAddress,
		transport:     nil,
	}
}

func (m *Multiplexer) Start(ctx context.Context) error {
	transport := transport.NewSocketServer(m.listenAddress)
	m.transport = transport
	m.transport.OnMessage(func(message json.RawMessage) {
		// TODO: Implement message handling
		fmt.Println("Received message:", string(message))
	})
	m.transport.OnError(func(err error) {
		// TODO: Implement error handling
		fmt.Println("Error:", err)
	})
	m.transport.OnClose(func() {
		fmt.Println("Close")
	})
	return m.transport.Start(ctx)
}
