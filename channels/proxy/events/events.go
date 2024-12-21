package events

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/protocol/mux"
)

type Events interface {
	EventMcpStarted()
	EventMcpResponseInitialize(initializeResponse *mcp.JsonRpcResponseInitializeResult)
	EventMcpToolsListResponse(toolsListResponse *mcp.JsonRpcResponseToolsListResult)

	EventMuxRequestToolCall(params *mux.JsonRpcRequestToolsCallParams, mcpReqId *jsonrpc.JsonRpcRequestId)

	EventMuxResponseProxyRegistered(registerResponse *mux.JsonRpcResponseProxyRegisterResult)
	EventMcpToolCallResponse(toolsCallResult *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, mcpReqId string)
}
