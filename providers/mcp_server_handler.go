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

	return &ProviderMcpServerHandler{
		logger:              logger,
		sdkServerDefinition: sdkServerDefinition,
		proxyRegistry:       proxyRegistry,
	}, nil
}

func (n *ProviderMcpServerHandler) ExecuteToolCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams) (interface{}, *jsonrpc.JsonRpcError) {
	n.logger.Info("OnToolCall", types.LogArg{"params": params})
	return nil, nil
}
