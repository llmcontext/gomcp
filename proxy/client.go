package proxy

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/llmcontext/gomcp/transport"
)

type Client struct {
	programName string
	args        []string
}

func NewClient(programName string, args []string) *Client {
	return &Client{
		programName: programName,
		args:        args,
	}
}

func (c *Client) Start() error {
	ctx, cancel := context.WithCancel(context.Background())

	// Use a wait group to wait for goroutines to complete
	var wg sync.WaitGroup

	// Listen for OS signals (e.g., Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGCHLD)

	proxyTransport := transport.NewStdioProxyClientTransport(c.programName, c.args)

	wg.Add(1)
	go func() {
		defer wg.Done()

		proxyClient := NewMCPProxyClient(proxyTransport)
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
