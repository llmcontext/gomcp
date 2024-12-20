package hub

import (
	"github.com/llmcontext/gomcp/channels/hub/events"
	"github.com/llmcontext/gomcp/channels/hubmcpserver"
	"github.com/llmcontext/gomcp/channels/hubmuxserver"
	"github.com/llmcontext/gomcp/jsonrpc"
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

func (s *StateManager) AsEvents() events.Events {
	return s
}

func (s *StateManager) EventNewProxyTools() {
	s.mcpServer.OnNewProxyTools()
}

func (s *StateManager) EventProxyToolCall(proxyId string, toolName string, toolArgs map[string]interface{}, id *jsonrpc.JsonRpcRequestId) {
	s.muxServer.OnProxyToolCall(proxyId, toolName, toolArgs, id)
}
