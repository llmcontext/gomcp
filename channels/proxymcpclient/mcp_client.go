package proxymcpclient

import (
	"context"

	"github.com/llmcontext/gomcp/channels/proxy/events"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type ProxyMcpClient struct {
	options *transport.ProxiedMcpServerDescription
	events  events.Events
	logger  types.Logger

	// context for proxy transport
	transport *transport.JsonRpcTransport
}

func NewProxyMcpClient(
	events events.Events,
	options *transport.ProxiedMcpServerDescription,
	logger types.Logger,
) *ProxyMcpClient {
	return &ProxyMcpClient{
		transport: nil,
		events:    events,
		options:   options,
		logger:    logger,
	}
}

func (c *ProxyMcpClient) Start(ctx context.Context) error {
	var err error
	errProxyChan := make(chan error, 1)

	// create the transport for the proxy client
	proxyTransport := transport.NewStdioProxyClientTransport(c.options)

	clientMcpJsonRpcTransport := transport.NewJsonRpcTransport(proxyTransport, "proxy - mcpclient", c.logger)
	c.transport = clientMcpJsonRpcTransport

	// First message to send is always an initialize request
	c.transport.OnStarted(func() {
		c.events.EventMcpStarted()
	})

	go func() {
		// start the proxy transport
		err = c.transport.Start(ctx, func(msg transport.JsonRpcMessage, jsonRpcTransport *transport.JsonRpcTransport) {
			c.logger.Debug("received message from proxy", types.LogArg{
				"message": msg,
				"method":  msg.Method,
				"name":    jsonRpcTransport.Name(),
			})
			c.handleIncomingMessage(msg)
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
	c.transport.Close()
}

func (c *ProxyMcpClient) SendNotification(method string) {
	notification := jsonrpc.JsonRpcRequest{
		JsonRpcVersion: jsonrpc.JsonRpcVersion,
		Method:         method,
	}
	c.transport.SendRequest(&notification, "")
}

func (s *ProxyMcpClient) SendJsonRpcResponse(response interface{}, id *jsonrpc.JsonRpcRequestId) {
	s.transport.SendResponse(&jsonrpc.JsonRpcResponse{
		JsonRpcVersion: jsonrpc.JsonRpcVersion,
		Id:             id,
		Result:         response,
		Error:          nil,
	})
}

func (s *ProxyMcpClient) SendResponse(response *jsonrpc.JsonRpcResponse) error {
	s.logger.Debug("JsonRpcResponse", types.LogArg{
		"response": response,
	})
	s.transport.SendResponse(response)
	return nil
}

func (s *ProxyMcpClient) SendRequestWithMethodAndParams(method string, params interface{}) {
	s.transport.SendRequestWithMethodAndParams(method, params, "")
}

func (s *ProxyMcpClient) SendError(code int, message string, id *jsonrpc.JsonRpcRequestId) {
	s.logger.Debug("JsonRpcError", types.LogArg{
		"code":    code,
		"message": message,
		"id":      id,
	})
	err := s.transport.SendError(code, message, id)
	if err != nil {
		s.logger.Error("failed to send error", types.LogArg{
			"error": err,
		})
	}
}
