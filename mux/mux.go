package mux

import (
	"context"
	"fmt"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/transport/socket"
	"github.com/llmcontext/gomcp/types"
)

type Multiplexer struct {
	listenAddress string
	socketServer  *socket.SocketServer
	sessions      []*MuxSession
	sessionCount  int
	logger        types.Logger
}

func NewMultiplexer(config *config.ProxyConfig, logger types.Logger) *Multiplexer {
	return &Multiplexer{
		listenAddress: config.ListenAddress,
		socketServer:  nil,
		sessions:      []*MuxSession{},
		sessionCount:  0,
		logger:        logger,
	}
}

func (m *Multiplexer) Start(ctx context.Context) error {
	// create transport
	m.socketServer = socket.NewSocketServer(m.listenAddress)

	m.socketServer.OnError(func(err error) {
		m.logger.Error("Error", types.LogArg{
			"error": err,
		})
	})

	m.socketServer.Start(ctx, func(transport types.Transport) {
		// we have a new session
		sessionId := fmt.Sprintf("s-%03d", m.sessionCount)
		m.sessionCount++
		m.logger.Info("new session", types.LogArg{
			"sessionId": sessionId,
		})
		subLogger := types.NewSubLogger(m.logger, types.LogArg{
			"sessionId": sessionId,
		})
		session := NewMuxSession(sessionId, transport, subLogger)
		m.sessions = append(m.sessions, session)
		// TODO: error group here?
		err := session.Start(ctx)
		if err != nil {
			m.logger.Error("Failed to start session", types.LogArg{
				"error": err,
			})
		}
	})

	return nil
}

func (m *Multiplexer) Close() {
	m.socketServer.Close()
	for _, session := range m.sessions {
		session.Close()
	}
}
