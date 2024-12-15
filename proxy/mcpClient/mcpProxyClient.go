package mcpClient

import (
	"context"

	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
	"github.com/llmcontext/gomcp/version"
)

type MCPProxyClientOptions struct {
	ProxyName               string
	CurrentWorkingDirectory string
	ProgramName             string
	ProgramArgs             []string
}

type MCPProxyClient struct {
	options MCPProxyClientOptions
	logger  types.Logger

	// context for mux transport
	muxJsonRpcTransport *transport.JsonRpcTransport
	// context for proxy transport
	proxyJsonRpcTransport *transport.JsonRpcTransport

	// serverInfo is the info about the server we are connected to
	serverInfo mcp.ServerInfo
	// tools is the list of tools available for the proxy
	tools []mcp.ToolDescription
}

func NewMCPProxyClient(
	proxyJsonRpcTransport *transport.JsonRpcTransport,
	muxJsonRpcTransport *transport.JsonRpcTransport,
	options MCPProxyClientOptions,
	logger types.Logger,
) *MCPProxyClient {
	return &MCPProxyClient{
		proxyJsonRpcTransport: proxyJsonRpcTransport,
		muxJsonRpcTransport:   muxJsonRpcTransport,
		options:               options,
		logger:                logger,
		serverInfo:            mcp.ServerInfo{},
		tools:                 []mcp.ToolDescription{},
	}
}

func (c *MCPProxyClient) Start(ctx context.Context) error {
	var err error

	err = c.proxyJsonRpcTransport.Start(ctx, func(msg transport.JsonRpcMessage) {
		c.logger.Debug("received message from proxy", types.LogArg{
			"message": msg,
		})
		c.handleMcpIncomingMessage(msg, c.proxyJsonRpcTransport)
	})
	if err != nil {
		return err
	}

	err = c.muxJsonRpcTransport.Start(ctx, func(msg transport.JsonRpcMessage) {
		c.logger.Debug("received message from mux", types.LogArg{
			"message": msg,
		})
		// c.handleMuxIncomingMessage(msg)
	})
	if err != nil {
		return err
	}

	// First message to send is always an initialize request

	// we create the parameters for the initialize request
	// the proxy does not have any capabilities
	params := mcp.JsonRpcRequestInitializeParams{
		ProtocolVersion: mcp.ProtocolVersion,
		Capabilities:    mcp.ClientCapabilities{},
		ClientInfo: mcp.ClientInfo{
			Name:    c.options.ProxyName,
			Version: version.Version,
		},
	}

	// we send the initialize request to the proxy
	err = c.proxyJsonRpcTransport.SendRequestWithMethodAndParams(
		mcp.RpcRequestMethodInitialize, params)
	if err != nil {
		c.logger.Error("failed to send initialize request", types.LogArg{
			"error": err,
		})
		return err
	}

	// Keep the main thread alive
	// will be interrupted by the context
	<-ctx.Done()

	// we close the transports
	c.proxyJsonRpcTransport.Close()
	c.muxJsonRpcTransport.Close()

	c.logger.Info("shutdown\n", types.LogArg{})

	return nil
}
