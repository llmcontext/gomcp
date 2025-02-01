package modelcontextprotocol

import (
	"context"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

type McpServerEventHandler interface {
	// tools
	ExecuteToolsList(ctx context.Context, logger types.Logger) (*mcp.JsonRpcResponseToolsListResult, *jsonrpc.JsonRpcError)
	ExecuteToolCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams, logger types.Logger) (types.ToolCallResult, *jsonrpc.JsonRpcError)

	// prompts
	ExecutePromptsList(ctx context.Context, logger types.Logger) (*mcp.JsonRpcResponsePromptsListResult, *jsonrpc.JsonRpcError)
	ExecutePromptGet(ctx context.Context, params *mcp.JsonRpcRequestPromptsGetParams, logger types.Logger) (types.PromptGetResult, *jsonrpc.JsonRpcError)
}
