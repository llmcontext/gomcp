package proxy

import (
	"context"
	"fmt"
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
	muxAddress              string
	currentWorkingDirectory string
	programName             string
	args                    []string
}

const (
	GomcpProxyClientName = "gomcp-proxy"
)

func NewProxyClient(muxAddress string, currentWorkingDirectory string, programName string, args []string) *ProxyClient {
	return &ProxyClient{
		muxAddress:              muxAddress,
		currentWorkingDirectory: currentWorkingDirectory,
		programName:             programName,
		args:                    args,
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

	// create a logger for proxy messages
	logger := NewProxyLogger()

	// start the mux client
	wg.Add(1)
	go func() {
		defer wg.Done()
		// create a transport for the mux client
		muxClientSocket := socket.NewSocketClient(c.muxAddress)

		// we try to start the mux client socket
		// let's get a transport for the mux client
		muxClientTransport, err := muxClientSocket.Start()
		if err != nil {
			logger.Error("error starting mux client socket", types.LogArg{"error": err})
			return
		}

		if muxClientTransport == nil {
			logger.Error("error starting mux client socket", types.LogArg{"error": err})
			return
		}

		muxJsonRpcTransport := transport.NewJsonRpcTransport(muxClientTransport, logger)

		// create the options for the proxy client
		options := mcpClient.MCPProxyClientOptions{
			ProxyName:               GomcpProxyClientName,
			CurrentWorkingDirectory: c.currentWorkingDirectory,
			ProgramName:             c.programName,
			ProgramArgs:             c.args,
		}

		// create the transport for the proxy client
		proxyTransport := transport.NewStdioProxyClientTransport(
			options.ProgramName,
			options.ProgramArgs,
		)

		proxyJsonRpcTransport := transport.NewJsonRpcTransport(proxyTransport, logger)

		// create the proxy client
		proxyClient := mcpClient.NewMCPProxyClient(
			proxyJsonRpcTransport,
			muxJsonRpcTransport,
			options,
			logger,
		)

		// start the proxy client and the queues
		proxyClient.Start(ctx)
	}()

	// Wait for a signal to stop the server
	sig := <-signalChan
	fmt.Fprintf(os.Stderr, "[proxy] Received an interrupt, shutting down... %s\n", sig)

	// Cancel the context to signal the goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Fprintf(os.Stderr, "[proxy] All goroutines have stopped. Exiting.\n")

	return nil
}
