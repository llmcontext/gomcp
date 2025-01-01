package modelcontextprotocol

import (
	"context"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
)

type McpClientNotifications interface {
	DoStopAfterListOfFeatures() bool
	OnServerInformation(serverName string, serverVersion string)
	OnToolsList(result *mcp.JsonRpcResponseToolsListResult, rpcError *jsonrpc.JsonRpcError)
	OnToolCallResponse(result *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, rpcError *jsonrpc.JsonRpcError)
}

type McpServerNotifications interface {
	OnToolCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams)
}
