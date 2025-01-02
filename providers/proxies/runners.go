package proxies

import (
	"context"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/modelcontextprotocol"
	"github.com/llmcontext/gomcp/modelcontextprotocol/mcpclient"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
	"github.com/llmcontext/gomcp/version"
	"golang.org/x/sync/errgroup"
)

type McpProxyRunner struct {
	proxyDefinition *ProxyDefinition
}

type McpRunners struct {
	proxies map[string]*McpProxyRunner
}

func NewMcpRunners() *McpRunners {
	return &McpRunners{
		proxies: make(map[string]*McpProxyRunner),
	}
}

func (r *McpRunners) AddProxy(proxy *ProxyDefinition) {
	r.proxies[proxy.ProxyId] = &McpProxyRunner{
		proxyDefinition: proxy,
	}
}

func (r *McpRunners) ExecuteToolCall(
	ctx context.Context,
	proxyDefinition *ProxyDefinition,
	params *mcp.JsonRpcRequestToolsCallParams,
	logger types.Logger,
) (types.ToolCallResult, *jsonrpc.JsonRpcError) {
	return nil, nil
}

func (t *McpProxyRunner) CreateMcpClient(ctx context.Context, logger types.Logger) {
	eg, egctx := errgroup.WithContext(ctx)

	// we need to create a MCP client that will execute the tool
	mcpClient := mcpclient.NewMcpClient(
		config.DefaultApplicationName,
		version.Version,
		t.AsMcpClientNotifications(),
		logger,
	)

	// we need to create a program description that will be used to start the proxy
	program := &transport.ProxiedMcpServerDescription{
		ProxyId:                 t.proxyDefinition.ProxyId,
		CurrentWorkingDirectory: t.proxyDefinition.WorkingDirectory,
		ProgramName:             t.proxyDefinition.ProgramName,
		ProgramArgs:             t.proxyDefinition.ProgramArguments,
	}

	eg.Go(func() error {
		err := mcpClient.StartWithMcpServer(egctx, program)
		if err != nil {
			logger.Error("error starting mux client", types.LogArg{"error": err})
		}

		return err
	})

	err := eg.Wait()
	if err != nil {
		logger.Error("error starting proxy client", types.LogArg{"error": err})
	}
	// TODO: we need to return the mcp client
}

func (t *McpProxyRunner) AsMcpClientNotifications() modelcontextprotocol.McpClientEventHandler {
	return t
}

func (t *McpProxyRunner) DoStopAfterListOfFeatures() bool {
	return false
}

func (t *McpProxyRunner) OnToolCallResponse(result *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, rpcError *jsonrpc.JsonRpcError) {

}

func (t *McpProxyRunner) OnServerInformation(serverName string, serverVersion string) {

}

func (t *McpProxyRunner) OnToolsList(result *mcp.JsonRpcResponseToolsListResult, rpcError *jsonrpc.JsonRpcError) {

}
