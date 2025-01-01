package providers

import (
	"context"

	"github.com/llmcontext/gomcp/modelcontextprotocol"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/providers/proxies"
	"github.com/llmcontext/gomcp/providers/sdk"
	"github.com/llmcontext/gomcp/types"
)

type ProviderMcpServerNotifications struct {
	logger              types.Logger
	sdkServerDefinition *sdk.SdkServerDefinition
	proxyRegistry       *proxies.ProxyRegistry
}

func NewProviderMcpServerNotifications(
	sdkServerDefinition *sdk.SdkServerDefinition,
	withProxies bool,
	logger types.Logger) (modelcontextprotocol.McpServerNotifications, error) {
	var proxyRegistry *proxies.ProxyRegistry
	var err error
	if withProxies {
		proxyRegistry, err = proxies.NewProxyRegistry()
		if err != nil {
			return nil, err
		}
	}

	return &ProviderMcpServerNotifications{
		logger:              logger,
		sdkServerDefinition: sdkServerDefinition,
		proxyRegistry:       proxyRegistry,
	}, nil
}

func (n *ProviderMcpServerNotifications) OnToolCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams) {
	n.logger.Info("OnToolCall", types.LogArg{"params": params})
}
