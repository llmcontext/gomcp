package proxies

import (
	"context"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

func (r *ProxyRegistry) ExecuteToolCall(
	ctx context.Context,
	params *mcp.JsonRpcRequestToolsCallParams,
	logger types.Logger,
) (types.ToolCallResult, *jsonrpc.JsonRpcError) {
	toolName := params.Name
	arguments := params.Arguments

	logger.Info("ExecuteToolCall", types.LogArg{
		"toolName":  toolName,
		"arguments": arguments,
	})

	proxy := r.GetProxy(toolName)
	if proxy == nil {
		return nil, &jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInternalError,
			Message: fmt.Sprintf("proxy %s not found", toolName),
		}
	}

	return nil, nil
}
