package hub

import (
	"github.com/llmcontext/gomcp/channels/hub/events"
	"github.com/llmcontext/gomcp/channels/hubmcpserver"
	"github.com/llmcontext/gomcp/channels/hubmuxserver"
	"github.com/llmcontext/gomcp/types"
)

type StateManager struct {
	logger    types.Logger
	mcpServer *hubmcpserver.MCPServer
	muxServer *hubmuxserver.MuxServer
}

func NewStateManager(logger types.Logger) *StateManager {
	return &StateManager{
		logger: logger,
	}
}

func (s *StateManager) SetMcpServer(server *hubmcpserver.MCPServer) {
	s.mcpServer = server
}

func (s *StateManager) SetMuxServer(server *hubmuxserver.MuxServer) {
	s.muxServer = server
}

func (s *StateManager) AsEventsProcessor() events.EventsProcessor {
	return s
}

func (s *StateManager) EventNewProxyTools() {
	s.mcpServer.OnNewProxyTools()
}