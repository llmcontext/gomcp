package proxy

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/llmcontext/gomcp/defaults"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/mcpclient"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
	"github.com/llmcontext/gomcp/version"
	"golang.org/x/sync/errgroup"
)

type Proxy struct {
	program         *transport.ProxiedMcpServerDescription
	logger          types.Logger
	proxyRegistry   *ProxyRegistry
	proxyDefinition *ProxyDefinition
}

func NewProxy(program *transport.ProxiedMcpServerDescription, logger types.Logger) *Proxy {
	return &Proxy{
		program: program,
		logger:  logger,
	}
}

func (p *Proxy) Start() error {
	var err error
	// create a context for the errgroup
	ctx := context.Background()
	eg, egctx := errgroup.WithContext(ctx)

	// we retrieve the proxy registry
	proxyRegistry, err := NewProxyRegistry()
	if err != nil {
		return err
	}
	p.proxyRegistry = proxyRegistry

	// we prepare the structure to register the proxy
	p.proxyDefinition = &ProxyDefinition{
		ProxyId:          p.program.ProxyId,
		WorkingDirectory: p.program.CurrentWorkingDirectory,
		ProgramName:      p.program.ProgramName,
		ProgramArguments: p.program.ProgramArgs,
		Tools:            []*ProxyToolDefinition{},
	}

	// goroutine to listen for OS signals
	// this will be used to stop the proxy client
	eg.Go(func() error {
		// Listen for OS signals (e.g., Ctrl+C)
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGCHLD)

		select {
		case <-egctx.Done():
			close(signalChan)
			return egctx.Err()
		case sig := <-signalChan:
			p.logger.Info("Received an interrupt, shutting down...", types.LogArg{"signal": sig})
			return fmt.Errorf("received an interrupt (%s)", sig.String())
		}
	})

	mcpClient := mcpclient.NewMcpClient(
		defaults.DefaultApplicationName,
		version.Version,
		p.logger,
		p,
	)

	eg.Go(func() error {
		err = mcpClient.StartWithMcpServer(egctx, p.program)
		if err != nil {
			p.logger.Error("error starting mux client", types.LogArg{"error": err})
		}

		return err
	})

	err = eg.Wait()
	if err != nil {
		p.logger.Error("error starting proxy client", types.LogArg{"error": err})
	}

	// we register the proxy
	p.logger.Info("Registering proxy", types.LogArg{"proxyDefinition": p.proxyDefinition})
	err = p.proxyRegistry.AddProxy(p.proxyDefinition)
	if err != nil {
		p.logger.Error("error registering proxy", types.LogArg{"error": err})
	}

	p.logger.Info("All goroutines have stopped. Exiting.", types.LogArg{})

	return nil

}

func (p *Proxy) DoStopAfterListOfFeatures() bool {
	return true
}

func (p *Proxy) OnServerInformation(serverName string, serverVersion string) {
	p.logger.Info("Server information received", types.LogArg{"serverName": serverName, "serverVersion": serverVersion})
	p.proxyDefinition.ProxyServerName = serverName
	p.proxyDefinition.ProxyServerVersion = serverVersion
}

func (p *Proxy) OnToolsList(result *mcp.JsonRpcResponseToolsListResult, rpcError *jsonrpc.JsonRpcError) {
	p.logger.Info("Tools list received", types.LogArg{"result": result, "rpcError": rpcError})
	for _, tool := range result.Tools {
		p.proxyDefinition.Tools = append(p.proxyDefinition.Tools, &ProxyToolDefinition{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}
}

// should be called in a proxy client
func (p *Proxy) OnToolCallResponse(result *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, rpcError *jsonrpc.JsonRpcError) {
	p.logger.Error("Tool call response received", types.LogArg{"result": result, "reqId": reqId, "rpcError": rpcError})
}
