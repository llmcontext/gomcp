package proxies

import (
	"context"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

func (r *ProxyRegistry) ExecuteToolCall(
	ctx context.Context,
	proxyDefinition *ProxyDefinition,
	proxyTool *ProxyToolDefinition,
	params *mcp.JsonRpcRequestToolsCallParams,
	logger types.Logger,
) (types.ToolCallResult, *jsonrpc.JsonRpcError) {
	// we delegate the execution to the mcp runners
	return r.mcpRunners.ExecuteToolCall(ctx, proxyDefinition, params, logger)
}
