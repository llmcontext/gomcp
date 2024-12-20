package events

import (
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/protocol/mux"
)

type Events interface {
	EventMcpStarted()
	EventMcpInitializeResponse(initializeResponse *mcp.JsonRpcResponseInitializeResult)
	EventMcpToolsListResponse(toolsListResponse *mcp.JsonRpcResponseToolsListResult)

	EventMuxProxyRegistered(registerResponse *mux.JsonRpcResponseProxyRegisterResult)
	EventMuxToolCall(name string, args map[string]interface{}, mcpReqId string)
}
