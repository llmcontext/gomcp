package hub

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/channels/hub/events"
	"github.com/llmcontext/gomcp/channels/hubmcpserver"
	"github.com/llmcontext/gomcp/channels/hubmuxserver"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

type ClientInfo struct {
	name    string
	version string
}

type StateManager struct {
	// mcp related state
	serverName          string
	serverVersion       string
	clientInfo          *ClientInfo
	isClientInitialized bool

	// mux related state

	logger    types.Logger
	mcpServer *hubmcpserver.MCPServer
	muxServer *hubmuxserver.MuxServer
}

func NewStateManager(serverName string, serverVersion string, logger types.Logger) *StateManager {
	return &StateManager{
		serverName:          serverName,
		serverVersion:       serverVersion,
		isClientInitialized: false,
		logger:              logger,
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

func (s *StateManager) EventMcpRequestInitialize(params *mcp.JsonRpcRequestInitializeParams, reqId *jsonrpc.JsonRpcRequestId) {
	// store client information
	if params.ProtocolVersion != mcp.ProtocolVersion {
		s.logger.Error("protocol version mismatch", types.LogArg{
			"expected": mcp.ProtocolVersion,
			"received": params.ProtocolVersion,
		})
	}
	s.clientInfo = &ClientInfo{
		name:    params.ClientInfo.Name,
		version: params.ClientInfo.Version,
	}

	// prepare response
	response := mcp.JsonRpcResponseInitializeResult{
		ProtocolVersion: mcp.ProtocolVersion,
		Capabilities: mcp.ServerCapabilities{
			Tools: &mcp.ServerCapabilitiesTools{
				ListChanged: jsonrpc.BoolPtr(false),
			},
			Prompts: &mcp.ServerCapabilitiesPrompts{
				ListChanged: jsonrpc.BoolPtr(false),
			},
		},
		ServerInfo: mcp.ServerInfo{Name: s.serverName, Version: s.serverVersion},
	}
	s.mcpServer.SendJsonRpcResponse(&response, reqId)

}

func (s *StateManager) EventMcpNotificationInitialized() {
	// that's a notification, no response is needed
	s.isClientInitialized = true
}

func (s *StateManager) EventNewProxyTools() {
	s.mcpServer.OnNewProxyTools()
}

func (s *StateManager) EventProxyToolCall(proxyId string, toolName string, toolArgs map[string]interface{}, mcpReqId string) {
	s.muxServer.OnProxyToolCall(proxyId, toolName, toolArgs, mcpReqId)
}

func (s *StateManager) EventMcpError(code int, message string, data *json.RawMessage, id *jsonrpc.JsonRpcRequestId) {
	s.mcpServer.SendError(code, message, id)
}
