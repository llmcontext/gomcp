package proxy

import (
	"github.com/llmcontext/gomcp/channels/proxy/events"
	"github.com/llmcontext/gomcp/channels/proxymcpclient"
	"github.com/llmcontext/gomcp/channels/proxymuxclient"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type StateManager struct {
	logger    types.Logger
	options   *transport.ProxiedMcpServerDescription
	muxClient *proxymuxclient.ProxyMuxClient
	mcpClient *proxymcpclient.ProxyMcpClient

	// serverInfo is the info about the MCP server we are connected to
	serverInfo mcp.ServerInfo
	proxyId    string
}

func NewStateManager(options *transport.ProxiedMcpServerDescription, logger types.Logger) *StateManager {
	return &StateManager{
		options:    options,
		logger:     logger,
		serverInfo: mcp.ServerInfo{},
		proxyId:    "",
	}
}

func (s *StateManager) Stop(err error) {
	s.logger.Info("stopping state manager", types.LogArg{"error": err})
}

func (s *StateManager) SetMuxClient(muxClient *proxymuxclient.ProxyMuxClient) {
	s.muxClient = muxClient
}

func (s *StateManager) SetProxyClient(proxyClient *proxymcpclient.ProxyMcpClient) {
	s.mcpClient = proxyClient
}

func (s *StateManager) AsEvents() events.Events {
	return s
}

func (s *StateManager) EventMcpStarted() {
	s.mcpClient.SendInitializeRequest()
}

func (s *StateManager) EventMcpInitializeResponse(resp *mcp.JsonRpcResponseInitializeResult) {
	s.logger.Info("event mcp initialize response", types.LogArg{
		"name":    resp.ServerInfo.Name,
		"version": resp.ServerInfo.Version,
	})

	// we update the server information
	s.serverInfo.Name = resp.ServerInfo.Name
	s.serverInfo.Version = resp.ServerInfo.Version

	// we can now report the tools list to the mux server
	s.muxClient.SendProxyRegistrationRequest(s.options, s.serverInfo)

}

func (s *StateManager) EventMcpToolsListResponse(resp *mcp.JsonRpcResponseToolsListResult) {
	s.logger.Info("event mcp tools list response", types.LogArg{
		"tools": resp.Tools,
	})

	// we send the "tools/register" request
	s.muxClient.SendToolsRegisterRequest(resp.Tools)
}

func (s *StateManager) EventMuxResponseProxyRegistered(registerResponse *mux.JsonRpcResponseProxyRegisterResult) {
	s.logger.Info("event mux proxy registered", types.LogArg{
		"sessionId":  registerResponse.SessionId,
		"proxyId":    registerResponse.ProxyId,
		"persistent": registerResponse.Persistent,
		"denied":     registerResponse.Denied,
	})

	// we store the proxy id
	s.proxyId = registerResponse.ProxyId

	// we send the "notifications/initialized" notification
	s.mcpClient.SendInitializedNotification()

	// we send the "tools/list" request
	s.mcpClient.SendToolsListRequest()

}

func (s *StateManager) EventMuxRequestToolCall(params *mux.JsonRpcRequestToolsCallParams, reqId *jsonrpc.JsonRpcRequestId) {
	s.logger.Info("EventMuxRequestToolCall", types.LogArg{
		"name":     params.Name,
		"args":     params.Args,
		"mcpReqId": params.McpReqId,
	})

	req := mcp.JsonRpcRequestToolsCallParams{
		Name:      params.Name,
		Arguments: params.Args,
	}

	// we forward the tool call to the mcp client
	s.mcpClient.SendRequestWithMethodAndParams(mcp.RpcRequestMethodToolsCall, req, reqId)
}

func (s *StateManager) EventMcpToolCallResponse(toolsCallResult *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, mcpReqId string) {
	s.logger.Info("event mcp tool call response", types.LogArg{
		"content": toolsCallResult.Content,
		"isError": toolsCallResult.IsError,
	})
	s.muxClient.SendToolCallResponse(toolsCallResult, reqId, mcpReqId)
}
