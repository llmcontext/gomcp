package hub

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/channels/hub/events"
	"github.com/llmcontext/gomcp/channels/hubmcpserver"
	"github.com/llmcontext/gomcp/channels/hubmuxserver"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/prompts"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/tools"
	"github.com/llmcontext/gomcp/types"
)

// MCP client information (eg Claude)
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
	toolsRegistry       *tools.ToolsRegistry
	promptsRegistry     *prompts.PromptsRegistry

	logger       types.Logger
	mcpServer    *hubmcpserver.MCPServer
	muxServer    *hubmuxserver.MuxServer
	reqIdMapping *jsonrpc.ReqIdMapping
}

func NewStateManager(
	serverName string,
	serverVersion string,
	toolsRegistry *tools.ToolsRegistry,
	promptsRegistry *prompts.PromptsRegistry,
	logger types.Logger,
) *StateManager {
	return &StateManager{
		serverName:          serverName,
		serverVersion:       serverVersion,
		isClientInitialized: false,
		toolsRegistry:       toolsRegistry,
		promptsRegistry:     promptsRegistry,
		logger:              logger,
		reqIdMapping:        jsonrpc.NewReqIdMapping(),
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
				ListChanged: jsonrpc.BoolPtr(true),
			},
			Prompts: &mcp.ServerCapabilitiesPrompts{
				ListChanged: jsonrpc.BoolPtr(true),
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

func (s *StateManager) EventMcpRequestToolsList(params *mcp.JsonRpcRequestToolsListParams, reqId *jsonrpc.JsonRpcRequestId) {
	// we query the tools registry
	tools := s.toolsRegistry.GetListOfTools()

	var response = mcp.JsonRpcResponseToolsListResult{
		Tools: make([]mcp.ToolDescription, 0, len(tools)),
	}

	// we build the response
	for _, tool := range tools {
		// schemaBytes, _ := json.Marshal(tool.InputSchema)
		response.Tools = append(response.Tools, mcp.ToolDescription{
			Name:        tool.ToolName,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}

	s.mcpServer.SendJsonRpcResponse(&response, reqId)
}

func (s *StateManager) EventMcpRequestToolsCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams, reqId *jsonrpc.JsonRpcRequestId) {
	// we get the tool name and arguments
	toolName := params.Name
	toolArgs := params.Arguments

	// let's check if the tool exists and is a proxy
	isProxy, proxyId, err := s.toolsRegistry.IsProxyTool(toolName)
	if err != nil {
		s.mcpServer.SendError(jsonrpc.RpcInternalError, fmt.Sprintf("tool not found: %v", err), reqId)
		return
	}

	// handle proxy tools
	if isProxy {
		session := s.muxServer.GetSessionByProxyId(proxyId)
		if session == nil {
			s.mcpServer.SendError(jsonrpc.RpcInternalError, "session not found", reqId)
			return
		}
		// we send the request to the proxy
		params := &mux.JsonRpcRequestToolsCallParams{
			Name: toolName,
			Args: toolArgs,
		}
		// we send the request to the proxy
		// we keep track of the request id for that tool call in the session extra parameters
		muxReqId, err := session.SendRequestWithMethodAndParams(mux.RpcRequestMethodCallTool, params)
		if err != nil {
			s.mcpServer.SendError(jsonrpc.RpcInternalError, fmt.Sprintf("failed to send request to proxy: %v", err), reqId)
			return
		}
		// we keep track of the mapping between the mcp request id
		// and the mux request id
		s.reqIdMapping.AddMapping(muxReqId, reqId)
	} else {
		// this is a direct tool call (SDK built-in tool)
		// let's call the tool
		response, err := s.toolsRegistry.CallTool(ctx, toolName, toolArgs)
		if err != nil {
			s.mcpServer.SendError(jsonrpc.RpcInternalError, fmt.Sprintf("tool call failed: %v", err), reqId)
			return
		}
		s.mcpServer.SendJsonRpcResponse(&response, reqId)
	}
}

func (s *StateManager) EventMcpRequestResourcesList(params *mcp.JsonRpcRequestResourcesListParams, reqId *jsonrpc.JsonRpcRequestId) {
	var response = mcp.JsonRpcResponseResourcesListResult{
		Resources: make([]mcp.ResourceDescription, 0),
	}

	s.mcpServer.SendJsonRpcResponse(&response, reqId)
}

func (s *StateManager) EventMcpRequestPromptsList(params *mcp.JsonRpcRequestPromptsListParams, reqId *jsonrpc.JsonRpcRequestId) {
	var response = mcp.JsonRpcResponsePromptsListResult{
		Prompts: make([]mcp.PromptDescription, 0),
	}

	prompts := s.promptsRegistry.GetListOfPrompts()
	for _, prompt := range prompts {
		arguments := make([]mcp.PromptArgumentDescription, 0, len(prompt.Arguments))
		for _, argument := range prompt.Arguments {
			arguments = append(arguments, mcp.PromptArgumentDescription{
				Name:        argument.Name,
				Description: argument.Description,
				Required:    argument.Required,
			})
		}
		response.Prompts = append(response.Prompts, mcp.PromptDescription{
			Name:        prompt.Name,
			Description: prompt.Description,
			Arguments:   arguments,
		})
	}

	s.mcpServer.SendJsonRpcResponse(&response, reqId)
}

func (s *StateManager) EventMcpRequestPromptsGet(params *mcp.JsonRpcRequestPromptsGetParams, reqId *jsonrpc.JsonRpcRequestId) {
	var templateArgs = map[string]string{}
	// copy the arguments, as strings
	for key, value := range params.Arguments {
		templateArgs[key] = fmt.Sprintf("%v", value)
	}
	promptName := params.Name

	response, err := s.promptsRegistry.GetPrompt(promptName, templateArgs)
	if err != nil {
		s.mcpServer.SendError(jsonrpc.RpcInvalidParams, fmt.Sprintf("prompt processing error: %s", err), reqId)
		return
	}

	// marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		s.mcpServer.SendError(jsonrpc.RpcInternalError, "failed to marshal response", reqId)
	}
	jsonResponse := json.RawMessage(responseBytes)

	// we send the response
	s.mcpServer.SendJsonRpcResponse(&jsonResponse, reqId)
}

func (s *StateManager) EventNewProxyTools() {
	// s.mcpServer.OnNewProxyTools()
}

func (s *StateManager) EventMcpError(code int, message string, data *json.RawMessage, id *jsonrpc.JsonRpcRequestId) {
	s.mcpServer.SendError(code, message, id)
}

func (s *StateManager) EventMuxRequestProxyRegister(proxyId string, params *mux.JsonRpcRequestProxyRegisterParams, reqId *jsonrpc.JsonRpcRequestId) {
	// we need to store the proxy id in the session
	session := s.muxServer.GetSessionByProxyId(proxyId)
	if session == nil {
		s.logger.Error("session not found", types.LogArg{
			"proxyId": proxyId,
		})
		return
	}
	session.SetSessionInformation(proxyId, params.ServerInfo.Name)

	// for now we accept all requests
	result := mux.JsonRpcResponseProxyRegisterResult{
		SessionId:  session.SessionId(),
		ProxyId:    proxyId,
		Persistent: params.Persistent,
		Denied:     false,
	}
	session.SendJsonRpcResponse(&result, reqId)
}

func (s *StateManager) EventMuxRequestToolsRegister(proxyId string, params *mux.JsonRpcRequestToolsRegisterParams, reqId *jsonrpc.JsonRpcRequestId) {
	// we need to store the proxy id in the session
	session := s.muxServer.GetSessionByProxyId(proxyId)
	if session == nil {
		s.logger.Error("session not found", types.LogArg{
			"proxyId": proxyId,
		})
		return
	}

	toolProvider, err := s.toolsRegistry.RegisterProxyToolProvider(proxyId, session.ProxyName())
	if err != nil {
		s.logger.Error("Failed to register proxy tool provider", types.LogArg{
			"error": err,
		})
		return
	}
	for _, tool := range params.Tools {
		err := toolProvider.AddProxyTool(tool.Name, tool.Description, tool.InputSchema)
		if err != nil {
			s.logger.Error("Failed to add proxy tool", types.LogArg{
				"error": err,
			})
			return
		}
	}

	// we need to prepare the tool provider so that it can be used by the hub
	err = s.toolsRegistry.PrepareProxyToolProvider(toolProvider)
	if err != nil {
		s.logger.Error("Failed to prepare proxy tool provider", types.LogArg{
			"error": err,
		})
		return
	}

	// send the notification to the MCP client
	// so that it will refresh the tools list
	s.mcpServer.SendNotification(mcp.RpcNotificationMethodToolsListChanged)

}

func (s *StateManager) EventMuxResponseToolCall(toolsCallResult *mux.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId) {
	// we need to find the mcp request id for the given mux request id
	mcpReqId := s.reqIdMapping.GetMapping(reqId)
	// we send the response to the mcp client
	s.logger.Info("EventMuxResponseToolCall", types.LogArg{
		"mcpReqId": mcpReqId,
		"reqId":    reqId,
		"result":   toolsCallResult,
	})
	mcpResponse := &mcp.JsonRpcResponseToolsCallResult{
		Content: toolsCallResult.Content,
		IsError: toolsCallResult.IsError,
	}
	s.mcpServer.SendJsonRpcResponse(mcpResponse, mcpReqId)
}
