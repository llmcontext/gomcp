package proxymuxclient

import (
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMuxClient) handleIncomingMessage(message transport.JsonRpcMessage) error {
	if message.Response != nil {
		c.logger.Error("received JsonRpcResponse", types.LogArg{
			"response": message.Response,
			"method":   message.Method,
		})
	} else if message.Request != nil {
		c.logger.Error("received JsonRpcRequest", types.LogArg{
			"request": message.Request,
			"method":  message.Method,
		})
	}
	return nil
}
