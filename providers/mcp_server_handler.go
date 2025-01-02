package providers

import (
	"context"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/modelcontextprotocol"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/providers/proxies"
	"github.com/llmcontext/gomcp/providers/sdk"
	"github.com/llmcontext/gomcp/types"
)

type ProviderMcpServerHandler struct {
	logger              types.Logger
	sdkServerDefinition *sdk.SdkServerDefinition
	proxyRegistry       *proxies.ProxyRegistry
}

func NewProviderMcpServerHandler(
	sdkServerDefinition *sdk.SdkServerDefinition,
	withProxies bool,
	logger types.Logger) (modelcontextprotocol.McpServerEventHandler, error) {
	var proxyRegistry *proxies.ProxyRegistry
	var err error
	if withProxies {
		proxyRegistry, err = proxies.NewProxyRegistry()
		if err != nil {
			return nil, err
		}
		// prepare the proxy registry so that we can use it
		proxyRegistry.Prepare()
	}

	// prepare the server so that we can use it
	sdkServerDefinition.Prepare()

	return &ProviderMcpServerHandler{
		logger:              logger,
		sdkServerDefinition: sdkServerDefinition,
		proxyRegistry:       proxyRegistry,
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

	// get the tools from the proxies
	proxiesTools := n.proxyRegistry.GetProxies()
	for _, proxy := range proxiesTools {
		for _, tool := range proxy.GetTools() {
			result.Tools = append(result.Tools, mcp.ToolDescription{
				Name:        tool.Name,
				Description: tool.Description,
				InputSchema: tool.InputSchema,
			})
		}
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

	// check if the tool is available in the proxy
	proxyTool, proxy := n.proxyRegistry.GetTool(toolName)
	if proxyTool != nil {
		return proxyTool.ExecuteToolCall(ctx, proxy, params, logger)
	}

	// check if the tool is available in the sdk
	tool := n.sdkServerDefinition.GetTool(toolName)
	if tool != nil {
		return n.sdkServerDefinition.ExecuteToolCall(ctx, params, logger)
	}

	// if the tool is not found in the proxy or the sdk, return an error
	return nil, &jsonrpc.JsonRpcError{
		Code:    jsonrpc.RpcInternalError,
		Message: fmt.Sprintf("Tool %s not found", toolName),
	}
}
