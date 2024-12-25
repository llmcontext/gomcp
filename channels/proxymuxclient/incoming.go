package proxymuxclient

import (
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMuxClient) handleIncomingMessage(message transport.JsonRpcMessage) error {
	if message.Response != nil {
		response := message.Response
		switch message.Method {
		case mux.RpcRequestMethodProxyRegister:
			{
				response, err := mux.ParseJsonRpcResponseProxyRegister(response)
				if err != nil {
					c.logger.Error("error in handleProxyRegisterResponse", types.LogArg{
						"error": err,
					})
					return err
				}

				c.events.EventMuxResponseProxyRegistered(response)
			}
		default:
			c.logger.Error("received message with unexpected method", types.LogArg{
				"method":   message.Method,
				"response": message.Response,
				"c":        "6thr",
			})
		}
	} else if message.Request != nil {
		request := message.Request
		switch message.Method {
		case mux.RpcRequestMethodCallTool:
			{
				params, err := mux.ParseJsonRpcRequestToolsCallParams(request)
				if err != nil {
					c.logger.Error("error in handleToolCall", types.LogArg{
						"error": err,
					})
					return err
				}
				c.events.EventMuxRequestToolCall(params, request.Id)
			}
		default:
			c.logger.Error("received message with unexpected method", types.LogArg{
				"method":  message.Method,
				"request": message.Request,
				"c":       "ty4t",
			})
		}
	}
	return nil
}
