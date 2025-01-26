package modelcontextprotocol

import (
	"context"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

type McpServerEventHandler interface {
	ExecuteToolCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams, logger types.Logger) (types.ToolCallResult, *jsonrpc.JsonRpcError)
	ExecuteToolsList(ctx context.Context, logger types.Logger) (*mcp.JsonRpcResponseToolsListResult, *jsonrpc.JsonRpcError)
}
