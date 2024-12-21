package hubmuxserver

import (
	"context"
	"fmt"
	"slices"

	"github.com/llmcontext/gomcp/channels/hub/events"
	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/tools"
	"github.com/llmcontext/gomcp/transport/socket"
	"github.com/llmcontext/gomcp/types"
)

type MuxServer struct {
	listenAddress string
	socketServer  *socket.SocketServer
	sessions      []*MuxSession
	sessionCount  int
	logger        types.Logger
	events        events.Events
	toolsRegistry *tools.ToolsRegistry
}

// server inside the mcp server in charge of multiplexing multiple proxy clients
func NewMuxServer(config *config.ProxyConfig, events events.Events, toolsRegistry *tools.ToolsRegistry, logger types.Logger) *MuxServer {
	return &MuxServer{
		listenAddress: config.ListenAddress,
		socketServer:  nil,
		sessions:      []*MuxSession{},
		sessionCount:  0,
		logger:        logger,
		events:        events,
		toolsRegistry: toolsRegistry,
	}
}

func (m *MuxServer) Start(ctx context.Context) error {
	// create socket server to listen for new proxy client connections
	m.socketServer = socket.NewSocketServer(m.listenAddress)

	m.socketServer.OnError(func(err error) {
		m.logger.Error("Error", types.LogArg{
			"error": err,
		})
	})

	// the parameter is a function that will be called when
	// a new connection is established with a proxy client
	m.socketServer.Start(ctx, func(transport types.Transport) {
		// we have a new session
		m.sessionCount++
		sessionId := fmt.Sprintf("s-%03d", m.sessionCount)
		m.logger.Info("new session", types.LogArg{
			"sessionId": sessionId,
		})
		subLogger := types.NewSubLogger(m.logger, types.LogArg{
			"sessionId": sessionId,
		})

		// create a new session
		session := NewMuxSession(sessionId, transport, subLogger, m.toolsRegistry, m.events)
		m.sessions = append(m.sessions, session)

		// start the session processing in a goroutine
		// this is to avoid blocking the main thread
		go func() {
			// start the session processing
			err := session.Start(ctx)
			if err != nil {
				m.logger.Error("mux session error - removing it", types.LogArg{
					"sessionId": sessionId,
					"error":     err,
				})
				session.Close()
				// if the session fails to start, we remove it from the list of sessions
				m.sessions = slices.DeleteFunc(m.sessions, func(s *MuxSession) bool {
					return s.SessionId() == sessionId
				})
			}

		}()
	})

	<-ctx.Done()
	m.logger.Info("mux server - context cancelled, closing", types.LogArg{})
	m.Close()
	return ctx.Err()

}

func (m *MuxServer) Close() {
	m.socketServer.Close()
	for _, session := range m.sessions {
		session.Close()
	}
}

func (m *MuxServer) GetSessionByProxyId(proxyId string) *MuxSession {
	for _, session := range m.sessions {
		if session.proxyId == proxyId {
			return session
		}
	}
	return nil
}
