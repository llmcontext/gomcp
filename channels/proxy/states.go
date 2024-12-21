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
	"github.com/llmcontext/gomcp/version"
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
	params := mcp.JsonRpcRequestInitializeParams{
		ProtocolVersion: mcp.ProtocolVersion,
		Capabilities:    mcp.ClientCapabilities{},
		ClientInfo: mcp.ClientInfo{
			Name:    s.options.ProxyName,
			Version: version.Version,
		},
	}
	s.mcpClient.SendRequestWithMethodAndParams(mcp.RpcRequestMethodInitialize, params)

}

func (s *StateManager) EventMcpResponseInitialize(resp *mcp.JsonRpcResponseInitializeResult) {
	s.logger.Info("event mcp initialize response", types.LogArg{
		"name":    resp.ServerInfo.Name,
		"version": resp.ServerInfo.Version,
	})

	// we update the server information
	s.serverInfo.Name = resp.ServerInfo.Name
	s.serverInfo.Version = resp.ServerInfo.Version

	params := mux.JsonRpcRequestProxyRegisterParams{
		ProtocolVersion: mux.MuxProtocolVersion,
		Proxy: mux.ProxyDescription{
			WorkingDirectory: s.options.CurrentWorkingDirectory,
			Command:          s.options.ProgramName,
			Args:             s.options.ProgramArgs,
		},
		ServerInfo: mux.ServerInfo{
			Name:    resp.ServerInfo.Name,
			Version: resp.ServerInfo.Version,
		},
	}

	// we can now report the tools list to the mux server
	s.muxClient.SendRequestWithMethodAndParams(mux.RpcRequestMethodProxyRegister, params)

}

func (s *StateManager) EventMcpResponseToolsList(resp *mcp.JsonRpcResponseToolsListResult) {
	s.logger.Info("event mcp tools list response", types.LogArg{
		"tools": resp.Tools,
	})

	// we send the "tools/register" request to the mux server
	toolsMux := make([]mux.ToolDescription, len(resp.Tools))
	for i, tool := range resp.Tools {
		toolsMux[i] = mux.ToolDescription{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		}
	}
	params := mux.JsonRpcRequestToolsRegisterParams{
		Tools: toolsMux,
	}
	s.muxClient.SendRequestWithMethodAndParams(mux.RpcRequestMethodToolsRegister, params)
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
	s.mcpClient.SendNotification(mcp.RpcNotificationMethodInitialized)

	// we send the "tools/list" request
	s.mcpClient.SendRequestWithMethodAndParams(
		mcp.RpcRequestMethodToolsList, mcp.JsonRpcRequestToolsListParams{})

}

// this is a tool call from the hub
func (s *StateManager) EventMuxRequestToolCall(params *mux.JsonRpcRequestToolsCallParams, reqId *jsonrpc.JsonRpcRequestId) {
	s.logger.Info("EventMuxRequestToolCall", types.LogArg{
		"name": params.Name,
		"args": params.Args,
	})
	// we keep track of the request id coming from the hub
	// TODO: state management
	// mcpReqId := jsonrpc.RequestIdToString(reqId)

	req := mcp.JsonRpcRequestToolsCallParams{
		Name:      params.Name,
		Arguments: params.Args,
	}

	// we forward the tool call to the mcp client
	// keeping track in extra parameter the mcp request id
	s.mcpClient.SendRequestWithMethodAndParams(mcp.RpcRequestMethodToolsCall, req)
}

// got the response for the tool call from the mcp client
func (s *StateManager) EventMcpResponseToolCall(toolsCallResult *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, mcpReqId string) {
	s.logger.Info("event mcp tool call response", types.LogArg{
		"content": toolsCallResult.Content,
		"isError": toolsCallResult.IsError,
	})
	//s.muxClient.SendToolCallResponse(toolsCallResult, reqId, mcpReqId)
	params := mux.JsonRpcResponseToolsCallResult{
		Content: toolsCallResult.Content,
		IsError: toolsCallResult.IsError,
	}
	// we parse the req id is the one coming from the hub
	// and we send the response to the hub with that id
	hubReqId := jsonrpc.ReqIdStringToId(mcpReqId)
	s.muxClient.SendJsonRpcResponse(params, hubReqId)

}
