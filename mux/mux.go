package mux

import "github.com/llmcontext/gomcp/config"

type Multiplexer struct {
	listenAddress string
}

func NewMultiplexer(config *config.ProxyConfig) *Multiplexer {
	return &Multiplexer{
		listenAddress: config.ListenAddress,
	}
}

func (m *Multiplexer) Start() error {
	return nil
}
