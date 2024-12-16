package proxy

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/llmcontext/gomcp/proxy/mcpClient"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/transport/socket"
	"github.com/llmcontext/gomcp/types"
)

type ProxyClient struct {
	proxyInformation ProxyInformation
	logger           types.Logger
}

type ProxyInformation struct {
	MuxAddress              string
	CurrentWorkingDirectory string
	ProgramName             string
	Args                    []string
}

const (
	GomcpProxyClientName = "gomcp-proxy"
)

func NewProxyClient(proxyInformation ProxyInformation, logger types.Logger) *ProxyClient {
	return &ProxyClient{
		proxyInformation: proxyInformation,
		logger:           logger,
	}
}

func (c *ProxyClient) Start() error {
	// Use a wait group to wait for goroutines to complete
	var wg sync.WaitGroup

	// Listen for OS signals (e.g., Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGCHLD)

	// start the mux client
	// create a transport for the mux client
	muxClientSocket := socket.NewSocketClient(c.proxyInformation.MuxAddress)

	// we try to start the mux client socket
	// let's get a transport for the mux client
	muxClientTransport, err := muxClientSocket.Start()
	if err != nil {
		c.logger.Error("error starting mux client socket", types.LogArg{"error": err})
		return err
	}

	if muxClientTransport == nil {
		c.logger.Error("error starting mux client socket", types.LogArg{"error": err})
		return err
	}

	muxJsonRpcTransport := transport.NewJsonRpcTransport(muxClientTransport, "proxy client - gomcp (mux)", c.logger)

	// create the options for the proxy client
	options := mcpClient.MCPProxyClientOptions{
		ProxyName:               GomcpProxyClientName,
		CurrentWorkingDirectory: c.proxyInformation.CurrentWorkingDirectory,
		ProgramName:             c.proxyInformation.ProgramName,
		ProgramArgs:             c.proxyInformation.Args,
	}

	// create the transport for the proxy client
	proxyTransport := transport.NewStdioProxyClientTransport(
		options.ProgramName,
		options.ProgramArgs,
	)

	proxyJsonRpcTransport := transport.NewJsonRpcTransport(proxyTransport, "proxy - client (mcp)", c.logger)

	// create the proxy client
	proxyClient := mcpClient.NewMCPProxyClient(
		proxyJsonRpcTransport,
		muxJsonRpcTransport,
		options,
		c.logger,
	)

	// prepare a cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// start the proxy client and the queues
	errProxyChan, err := proxyClient.Start(ctx)
	if err != nil {
		c.logger.Error("error starting proxy client", types.LogArg{"error": err})
		cancel()
		return err
	}
	// Wait for either a signal, error, or context cancellation
	select {
	case sig := <-signalChan:
		c.logger.Info("Received an interrupt, shutting down...", types.LogArg{"signal": sig})
	case err := <-errProxyChan:
		c.logger.Error("Goroutine stopped with error", types.LogArg{"error": err})
	}

	// Cancel the context to signal the goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
	c.logger.Info("All goroutines have stopped. Exiting.", types.LogArg{})

	return nil
}
