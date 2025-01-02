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

func (t *ProxyToolDefinition) ExecuteToolCall(
	ctx context.Context,
	proxyDefinition *ProxyDefinition,
	params *mcp.JsonRpcRequestToolsCallParams,
	logger types.Logger,
) (types.ToolCallResult, *jsonrpc.JsonRpcError) {
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
		ProxyId:                 proxyDefinition.ProxyId,
		CurrentWorkingDirectory: proxyDefinition.WorkingDirectory,
		ProgramName:             proxyDefinition.ProgramName,
		ProgramArgs:             proxyDefinition.ProgramArguments,
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

	return nil, nil
}

func (t *ProxyToolDefinition) DoStopAfterListOfFeatures() bool {
	return false
}

func (t *ProxyToolDefinition) AsMcpClientNotifications() modelcontextprotocol.McpClientEventHandler {
	return t
}

func (t *ProxyToolDefinition) OnToolCallResponse(result *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, rpcError *jsonrpc.JsonRpcError) {

}

func (t *ProxyToolDefinition) OnServerInformation(serverName string, serverVersion string) {

}

func (t *ProxyToolDefinition) OnToolsList(result *mcp.JsonRpcResponseToolsListResult, rpcError *jsonrpc.JsonRpcError) {

}
