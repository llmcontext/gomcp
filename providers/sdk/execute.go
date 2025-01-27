package sdk

import (
	"context"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/providers/results"
	"github.com/llmcontext/gomcp/types"
)

func (n *SdkServerDefinition) ExecuteToolCall(
	ctx context.Context,
	params *mcp.JsonRpcRequestToolsCallParams,
	logger types.Logger,
) (types.ToolCallResult, *jsonrpc.JsonRpcError) {
	toolName := params.Name
	arguments := params.Arguments
	tool := n.GetTool(toolName)
	if tool == nil {
		return nil, &jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInternalError,
			Message: fmt.Sprintf("tool %s not found", toolName),
		}
	}

	// check if we already have a context for this tool
	if tool.toolContext == nil {
		// if not, we initialize the tool context
		err := n.serverInitFunction(ctx, logger)
		if err != nil {
			logger.Error("error initializing tool context", types.LogArg{
				"toolName": toolName,
				"error":    err,
			})
			return nil, &jsonrpc.JsonRpcError{
				Code:    jsonrpc.RpcInternalError,
				Message: fmt.Sprintf("tool %s - error initializing tool context: %v", toolName, err),
			}
		}
	}

	// let's create the output
	output := results.NewToolCallResult()

	errChan := make(chan *jsonrpc.JsonRpcError, 1)
	go func() {
		tool.toolProcessFunction(ctx, arguments, output, logger, errChan)
	}()

	// wait on context and errChan
	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		} else {
			return output, nil
		}
	case <-ctx.Done():
		return nil, &jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInternalError,
			Message: ctx.Err().Error(),
		}
	}
}
