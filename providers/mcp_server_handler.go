package providers

import (
	"context"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/modelcontextprotocol"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/providers/sdk"
	"github.com/llmcontext/gomcp/types"
)

type ProviderMcpServerHandler struct {
	logger              types.Logger
	sdkServerDefinition *sdk.SdkServerDefinition
}

func NewProviderMcpServerHandler(
	sdkServerDefinition *sdk.SdkServerDefinition,
	logger types.Logger) (modelcontextprotocol.McpServerEventHandler, error) {
	// prepare the server so that we can use it
	err := sdkServerDefinition.Prepare()
	if err != nil {
		return nil, err
	}

	return &ProviderMcpServerHandler{
		logger:              logger,
		sdkServerDefinition: sdkServerDefinition,
	}, nil
}

func (n *ProviderMcpServerHandler) ExecuteToolsList(ctx context.Context, logger types.Logger) (*mcp.JsonRpcResponseToolsListResult, *jsonrpc.JsonRpcError) {
	result := &mcp.JsonRpcResponseToolsListResult{
		Tools: make([]mcp.ToolDescription, 0, 10),
	}

	// get the tools from the sdk
	tools := n.sdkServerDefinition.GetListOfTools()
	for _, tool := range tools {
		result.Tools = append(result.Tools, mcp.ToolDescription{
			Name:        tool.ToolName,
			Description: tool.ToolDescription,
			InputSchema: tool.InputSchema,
		})
	}

	return result, nil
}

func (n *ProviderMcpServerHandler) ExecuteToolCall(
	ctx context.Context,
	params *mcp.JsonRpcRequestToolsCallParams,
	logger types.Logger,
) (types.ToolCallResult, *jsonrpc.JsonRpcError) {
	toolName := params.Name
	n.logger.Info("OnToolCall", types.LogArg{
		"toolName": toolName,
		"params":   params,
	})

	// check if the tool is available in the sdk
	tool := n.sdkServerDefinition.GetTool(toolName)
	if tool != nil {
		return n.sdkServerDefinition.ExecuteToolCall(ctx, params, logger)
	}

	// if the tool is not found in the sdk, return an error
	return nil, &jsonrpc.JsonRpcError{
		Code:    jsonrpc.RpcInternalError,
		Message: fmt.Sprintf("Tool %s not found", toolName),
	}
}

func (n *ProviderMcpServerHandler) ExecutePromptsList(ctx context.Context, logger types.Logger) (*mcp.JsonRpcResponsePromptsListResult, *jsonrpc.JsonRpcError) {
	var response = mcp.JsonRpcResponsePromptsListResult{
		Prompts: make([]mcp.PromptDescription, 0),
	}

	prompts := n.sdkServerDefinition.GetListOfPrompts()
	for _, prompt := range prompts {
		arguments := make([]mcp.PromptArgumentDescription, 0, len(prompt.Arguments))
		for _, argument := range prompt.Arguments {
			arguments = append(arguments, mcp.PromptArgumentDescription{
				Name:        argument.Name,
				Description: argument.Description,
				Required:    argument.Required,
			})
		}
		response.Prompts = append(response.Prompts, mcp.PromptDescription{
			Name:        prompt.Name,
			Description: prompt.Description,
			Arguments:   arguments,
		})
	}

	return &response, nil
}

func (n *ProviderMcpServerHandler) ExecutePromptGet(ctx context.Context, params *mcp.JsonRpcRequestPromptsGetParams, logger types.Logger) (types.PromptGetResult, *jsonrpc.JsonRpcError) {
	promptName := params.Name
	templateArgs := params.Arguments

	response, err := n.sdkServerDefinition.GetPrompt(promptName, templateArgs)
	if err != nil {
		return nil, &jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidParams,
			Message: fmt.Sprintf("prompt processing error: %s", err),
		}
	}

	return response, nil
}
