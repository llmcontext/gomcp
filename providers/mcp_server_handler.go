package providers

import (
	"context"

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
	}

	// prepare the server lifecycle methods
	sdkServerDefinition.PrepareLifecyles()

	return &ProviderMcpServerHandler{
		logger:              logger,
		sdkServerDefinition: sdkServerDefinition,
		proxyRegistry:       proxyRegistry,
	}, nil
}

func (n *ProviderMcpServerHandler) ExecuteToolCall(
	ctx context.Context,
	toolName string,
	params *mcp.JsonRpcRequestToolsCallParams,
	logger types.Logger,
) (types.ToolCallResult, *jsonrpc.JsonRpcError) {
	n.logger.Info("OnToolCall", types.LogArg{
		"toolName": toolName,
		"params":   params,
	})
	return nil, nil
}

func (n *ProviderMcpServerHandler) ExecuteToolsList(ctx context.Context, logger types.Logger) (*mcp.JsonRpcResponseToolsListResult, *jsonrpc.JsonRpcError) {
	result := &mcp.JsonRpcResponseToolsListResult{
		Tools: make([]mcp.ToolDescription, 0, 10),
	}

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
