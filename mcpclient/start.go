package mcpclient

import (
	"context"

	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (m *McpClient) StartWithMcpServer(program *transport.ProxiedMcpServerDescription) error {
	var err error
	m.logger.Info("Starting MCP client", types.LogArg{"program": program})
	errProxyChan := make(chan error, 1)

	// create a context that will be used to cancel the server and the inspector
	ctx := context.Background()

	// create the transport
	proxyTransport := transport.NewStdioProxyClientTransport(program)

	// create the json rpc transport
	jsonRpcTransport := transport.NewJsonRpcTransport(proxyTransport, "proxy - mcpclient", m.logger)
	m.jsonRpcTransport = jsonRpcTransport

	// we report that the MCP server is started
	jsonRpcTransport.OnStarted(func() {
		m.EventMcpStarted()
	})

	go func() {
		// start the proxy transport
		err = jsonRpcTransport.Start(ctx, func(msg transport.JsonRpcMessage, jsonRpcTransport *transport.JsonRpcTransport) {
			m.logger.Debug("received message from proxy", types.LogArg{
				"message": msg,
				"method":  msg.Method,
				"name":    jsonRpcTransport.Name(),
			})
			m.handleIncomingMessage(msg)
		})
		if err != nil {
			m.logger.Error("failed to start proxy transport", types.LogArg{
				"error": err,
			})
			errProxyChan <- err
		}
	}()

	select {
	case err := <-errProxyChan:
		m.Close()
		return err
	case <-ctx.Done():
		m.Close()
		return ctx.Err()
	}
}

func (m *McpClient) Close() {
	m.jsonRpcTransport.Close()
}
