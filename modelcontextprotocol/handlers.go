package modelcontextprotocol

import (
	"context"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
)

type McpClientEventHandler interface {
	DoStopAfterListOfFeatures() bool
	OnServerInformation(serverName string, serverVersion string)
	OnToolsList(result *mcp.JsonRpcResponseToolsListResult, rpcError *jsonrpc.JsonRpcError)
	OnToolCallResponse(result *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, rpcError *jsonrpc.JsonRpcError)
}

type McpServerEventHandler interface {
	ExecuteToolCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams) (interface{}, *jsonrpc.JsonRpcError)
}
