package proxy

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/llmcontext/gomcp/endpoints/proxymcpclient"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/transport/socket"
	"github.com/llmcontext/gomcp/types"
	"golang.org/x/sync/errgroup"
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
	// create a context for the errgroup
	ctx := context.Background()
	eg, egctx := errgroup.WithContext(ctx)

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

	eg.Go(func() error {
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
			return fmt.Errorf("error creating mux transport")
		}

		// create the json rpc transport for the mux client
		muxJsonRpcTransport := transport.NewJsonRpcTransport(muxClientTransport, "proxy client - gomcp (mux)", c.logger)

		// create the options for the proxy client
		options := proxymcpclient.MCPProxyClientOptions{
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
		proxyClient := proxymcpclient.NewMCPProxyClient(
			proxyJsonRpcTransport,
			muxJsonRpcTransport,
			options,
			c.logger,
		)

		err = proxyClient.Start(egctx)
		if err != nil {
			c.logger.Error("error starting proxy client", types.LogArg{"error": err})
		}
		return err
	})

	err := eg.Wait()
	if err != nil {
		c.logger.Error("error starting proxy client", types.LogArg{"error": err})
	}

	c.logger.Info("All goroutines have stopped. Exiting.", types.LogArg{})

	return nil
}
