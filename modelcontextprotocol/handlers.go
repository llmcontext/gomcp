package modelcontextprotocol

import (
	"context"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

type McpClientEventHandler interface {
	DoStopAfterListOfFeatures() bool
	OnServerInformation(serverName string, serverVersion string)
	OnToolsList(result *mcp.JsonRpcResponseToolsListResult, rpcError *jsonrpc.JsonRpcError)
	OnToolCallResponse(result *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, rpcError *jsonrpc.JsonRpcError)
}

type McpServerEventHandler interface {
	ExecuteToolCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams, logger types.Logger) (types.ToolCallResult, *jsonrpc.JsonRpcError)
	ExecuteToolsList(ctx context.Context, logger types.Logger) (*mcp.JsonRpcResponseToolsListResult, *jsonrpc.JsonRpcError)
}
