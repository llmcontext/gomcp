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
	// prepare a cancelable context
	ctx, cancel := context.WithCancel(context.Background())

	// Use a wait group to wait for goroutines to complete
	var wg sync.WaitGroup

	// Listen for OS signals (e.g., Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGCHLD)

	// start the mux client
	wg.Add(1)
	go func() {
		defer wg.Done()
		// create a transport for the mux client
		muxClientSocket := socket.NewSocketClient(c.proxyInformation.MuxAddress)

		// we try to start the mux client socket
		// let's get a transport for the mux client
		muxClientTransport, err := muxClientSocket.Start()
		if err != nil {
			c.logger.Error("error starting mux client socket", types.LogArg{"error": err})
			return
		}

		if muxClientTransport == nil {
			c.logger.Error("error starting mux client socket", types.LogArg{"error": err})
			return
		}

		muxJsonRpcTransport := transport.NewJsonRpcTransport(muxClientTransport, c.logger)

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

		proxyJsonRpcTransport := transport.NewJsonRpcTransport(proxyTransport, c.logger)

		// create the proxy client
		proxyClient := mcpClient.NewMCPProxyClient(
			proxyJsonRpcTransport,
			muxJsonRpcTransport,
			options,
			c.logger,
		)

		// start the proxy client and the queues
		proxyClient.Start(ctx)
	}()

	// Wait for a signal to stop the server
	sig := <-signalChan
	c.logger.Info("Received an interrupt, shutting down...", types.LogArg{"signal": sig})

	// Cancel the context to signal the goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
	c.logger.Info("All goroutines have stopped. Exiting.", types.LogArg{})

	return nil
}
