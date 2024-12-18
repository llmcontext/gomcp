package proxymcpclient

import (
	"context"

	"github.com/llmcontext/gomcp/channels"
	"github.com/llmcontext/gomcp/channels/proxy/events"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
	"github.com/llmcontext/gomcp/version"
)

type ProxyMcpClient struct {
	options *channels.ProxiedMcpServerDescription
	events  *events.Events
	logger  types.Logger

	// context for proxy transport
	proxyJsonRpcTransport *transport.JsonRpcTransport

	// tools is the list of tools available for the proxy
	tools []mcp.ToolDescription
}

func NewProxyMcpClient(
	proxyJsonRpcTransport *transport.JsonRpcTransport,
	events *events.Events,
	options *channels.ProxiedMcpServerDescription,
	logger types.Logger,
) *ProxyMcpClient {
	return &ProxyMcpClient{
		proxyJsonRpcTransport: proxyJsonRpcTransport,
		events:                events,
		options:               options,
		logger:                logger,
		//serverInfo:            mcp.ServerInfo{},
		tools: []mcp.ToolDescription{},
	}
}

func (c *ProxyMcpClient) Start(ctx context.Context) error {
	var err error
	errProxyChan := make(chan error, 1)

	// First message to send is always an initialize request
	c.proxyJsonRpcTransport.OnStarted(func() {
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
		}
	})

	go func() {
		// start the proxy transport
		err = c.proxyJsonRpcTransport.Start(ctx, func(msg transport.JsonRpcMessage, jsonRpcTransport *transport.JsonRpcTransport) {
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
	c.proxyJsonRpcTransport.Close()
}

func (c *ProxyMcpClient) SendInitializedNotification() {
	notification := jsonrpc.NewJsonRpcNotification(mcp.RpcNotificationMethodInitialized)
	c.proxyJsonRpcTransport.SendRequest(notification)
}

func (c *ProxyMcpClient) SendToolsListRequest() {
	params := mcp.JsonRpcRequestToolsListParams{}
	c.proxyJsonRpcTransport.SendRequestWithMethodAndParams(
		mcp.RpcRequestMethodToolsList, params)
}
