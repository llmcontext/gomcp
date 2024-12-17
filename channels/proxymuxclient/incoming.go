package proxymuxclient

import (
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMuxClient) handleIncomingMessage(message transport.JsonRpcMessage) error {
	if message.Response != nil {
		c.logger.Error("received JsonRpcResponse", types.LogArg{
			"response": message.Response,
			"method":   message.Method,
		})
		switch message.Method {
		case mux.RpcRequestMethodProxyRegister:
			c.handleProxyRegisterResponse(message.Response)
		default:
			c.logger.Error("received message with unexpected method", types.LogArg{
				"method":   message.Method,
				"response": message.Response,
			})
		}
	} else if message.Request != nil {
		c.logger.Error("received JsonRpcRequest", types.LogArg{
			"request": message.Request,
			"method":  message.Method,
		})
	}
	return nil
}
