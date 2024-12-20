package proxymcpclient

import (
	"context"

	"github.com/llmcontext/gomcp/channels/proxy/events"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
	"github.com/llmcontext/gomcp/version"
)

type ProxyMcpClient struct {
	options *transport.ProxiedMcpServerDescription
	events  events.Events
	logger  types.Logger

	// context for proxy transport
	clientMcpJsonRpcTransport *transport.JsonRpcTransport

	// tools is the list of tools available for the proxy
	tools []mcp.ToolDescription
}

func NewProxyMcpClient(
	events events.Events,
	options *transport.ProxiedMcpServerDescription,
	logger types.Logger,
) *ProxyMcpClient {
	return &ProxyMcpClient{
		clientMcpJsonRpcTransport: nil,
		events:                    events,
		options:                   options,
		logger:                    logger,
		//serverInfo:            mcp.ServerInfo{},
		tools: []mcp.ToolDescription{},
	}
}

func (c *ProxyMcpClient) Start(ctx context.Context) error {
	var err error
	errProxyChan := make(chan error, 1)

	// create the transport for the proxy client
	proxyTransport := transport.NewStdioProxyClientTransport(c.options)

	clientMcpJsonRpcTransport := transport.NewJsonRpcTransport(proxyTransport, "proxy - mcpclient", c.logger)
	c.clientMcpJsonRpcTransport = clientMcpJsonRpcTransport

	// First message to send is always an initialize request
	c.clientMcpJsonRpcTransport.OnStarted(func() {
		c.events.EventMcpStarted()
	})

	go func() {
		// start the proxy transport
		err = c.clientMcpJsonRpcTransport.Start(ctx, func(msg transport.JsonRpcMessage, jsonRpcTransport *transport.JsonRpcTransport) {
			c.logger.Debug("received message from proxy", types.LogArg{
				"message": msg,
				"method":  msg.Method,
				"name":    jsonRpcTransport.Name(),
			})
			c.handleMcpIncomingMessage(msg)
		})
		if err != nil {
			c.logger.Error("failed to start proxy transport", types.LogArg{
				"error": err,
			})
			errProxyChan <- err
		}
	}()

	select {
	case err := <-errProxyChan:
		c.Close()
		return err
	case <-ctx.Done():
		c.Close()
		return ctx.Err()
	}
}

func (c *ProxyMcpClient) Close() {
	c.clientMcpJsonRpcTransport.Close()
}

func (c *ProxyMcpClient) SendInitializeRequest() {
	params := mcp.JsonRpcRequestInitializeParams{
		ProtocolVersion: mcp.ProtocolVersion,
		Capabilities:    mcp.ClientCapabilities{},
		ClientInfo: mcp.ClientInfo{
			Name:    c.options.ProxyName,
			Version: version.Version,
		},
	}
	c.clientMcpJsonRpcTransport.SendRequestWithMethodAndParams(mcp.RpcRequestMethodInitialize, params)
}

func (c *ProxyMcpClient) SendInitializedNotification() {
	notification := jsonrpc.NewJsonRpcNotification(mcp.RpcNotificationMethodInitialized)
	c.clientMcpJsonRpcTransport.SendRequest(notification)
}

func (c *ProxyMcpClient) SendToolsListRequest() {
	params := mcp.JsonRpcRequestToolsListParams{}
	c.clientMcpJsonRpcTransport.SendRequestWithMethodAndParams(
		mcp.RpcRequestMethodToolsList, params)
}
