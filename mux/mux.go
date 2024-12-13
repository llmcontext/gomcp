package mux

import (
	"context"
	"encoding/json"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/transport"
)

type Multiplexer struct {
	listenAddress string
	transport     *transport.SocketTransport
}

func NewMultiplexer(config *config.ProxyConfig) *Multiplexer {
	return &Multiplexer{
		listenAddress: config.ListenAddress,
		transport:     nil,
	}
}

func (m *Multiplexer) Start(ctx context.Context) error {
	m.transport = transport.NewSocketTransport(m.listenAddress)
	m.transport.OnMessage(func(message json.RawMessage) {
		// TODO: Implement message handling
	})
	m.transport.OnError(func(err error) {
		// TODO: Implement error handling
	})
	m.transport.OnClose(func() {
		// TODO: Implement close handling
	})
	return m.transport.Start(ctx)
}
