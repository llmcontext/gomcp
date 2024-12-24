package events

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/protocol/mux"
)

type Events interface {
	EventMcpStarted()
	EventMcpResponseInitialize(initializeResponse *mcp.JsonRpcResponseInitializeResult)
	EventMcpResponseToolsList(toolsListResponse *mcp.JsonRpcResponseToolsListResult)
	EventMcpResponseToolCall(toolsCallResult *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId)
	EventMcpResponseToolCallError(error *jsonrpc.JsonRpcError, reqId *jsonrpc.JsonRpcRequestId)

	EventMuxStarted()
	EventMuxRequestToolCall(params *mux.JsonRpcRequestToolsCallParams, mcpReqId *jsonrpc.JsonRpcRequestId)

	EventMuxResponseProxyRegistered(registerResponse *mux.JsonRpcResponseProxyRegisterResult)
}
