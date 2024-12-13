package proxy

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/llmcontext/gomcp/proxy/mcpClient"
	"github.com/llmcontext/gomcp/proxy/muxClient"
	"github.com/llmcontext/gomcp/transport"
)

type ProxyClient struct {
	muxAddress  string
	programName string
	args        []string
	muxClient   *muxClient.MuxClient
}

const (
	GomcpProxyClientName = "gomcp-proxy"
)

func NewProxyClient(muxAddress string, programName string, args []string) *ProxyClient {
	return &ProxyClient{
		muxAddress:  muxAddress,
		programName: programName,
		args:        args,
		muxClient:   nil,
	}
}

func (c *ProxyClient) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	// Use a wait group to wait for goroutines to complete
	var wg sync.WaitGroup

	// Listen for OS signals (e.g., Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGCHLD)

	// create a logger for proy messages
	logger := NewProxyLogger()
	c.muxClient = muxClient.NewMuxClient(c.muxAddress, logger)

	// create a mux client
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.muxClient.Start(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		proxyTransport := transport.NewStdioProxyClientTransport(c.programName, c.args)
		proxyClient := mcpClient.NewMCPProxyClient(proxyTransport, c.muxClient, GomcpProxyClientName, logger)
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
