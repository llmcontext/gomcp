package proxy

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/llmcontext/gomcp/channels/proxy/events"
	"github.com/llmcontext/gomcp/channels/proxymcpclient"
	"github.com/llmcontext/gomcp/channels/proxymuxclient"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
	"golang.org/x/sync/errgroup"
)

type ProxyClient struct {
	proxyInformation ProxyInformation
	logger           types.Logger
	stateManager     *StateManager
	events           events.Events
	options          *transport.ProxiedMcpServerDescription
}

type ProxyInformation struct {
	ProxyId                 string
	MuxAddress              string
	CurrentWorkingDirectory string
	ProgramName             string
	Args                    []string
}

const (
	GomcpProxyClientName = "gomcp-proxy"
)

func NewProxyClient(proxyInformation ProxyInformation, debug bool, logger types.Logger) *ProxyClient {
	options := transport.ProxiedMcpServerDescription{
		ProxyName:               GomcpProxyClientName,
		CurrentWorkingDirectory: proxyInformation.CurrentWorkingDirectory,
		ProgramName:             proxyInformation.ProgramName,
		ProgramArgs:             proxyInformation.Args,
		ProxyId:                 proxyInformation.ProxyId,
	}
	stateManager := NewStateManager(&options, logger)
	events := stateManager.AsEvents()

	return &ProxyClient{
		proxyInformation: proxyInformation,
		logger:           logger,
		options:          &options,
		stateManager:     stateManager,
		events:           events,
	}
}

func (c *ProxyClient) Start() error {
	var err error
	// create a context for the errgroup
	ctx := context.Background()
	eg, egctx := errgroup.WithContext(ctx)

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
			c.logger.Info("Received an interrupt, shutting down...", types.LogArg{"signal": sig})
			return fmt.Errorf("received an interrupt (%s)", sig.String())
		}
	})

	muxClient := proxymuxclient.NewProxyMuxClient(
		c.proxyInformation.MuxAddress,
		c.events,
		c.logger,
	)
	c.stateManager.SetMuxClient(muxClient)

	eg.Go(func() error {
		err = muxClient.Start(egctx)
		if err != nil {
			c.logger.Error("error starting mux client", types.LogArg{"error": err})
		}

		return err
	})

	// create the options for the proxy client

	// go routine for proxy mcp client
	eg.Go(func() error {
		// create the proxy client
		proxyClient := proxymcpclient.NewProxyMcpClient(
			c.events,
			c.options,
			c.logger,
		)

		c.stateManager.SetProxyClient(proxyClient)

		err := proxyClient.Start(egctx)
		if err != nil {
			c.logger.Error("error starting proxy client", types.LogArg{"error": err})
		}
		return err
	})

	err = eg.Wait()
	c.stateManager.Stop(err)
	if err != nil {
		c.logger.Error("error starting proxy client", types.LogArg{"error": err})
	}

	c.logger.Info("All goroutines have stopped. Exiting.", types.LogArg{})

	return nil
}
