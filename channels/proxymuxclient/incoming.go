package proxymuxclient

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMuxClient) handleIncomingMessage(message transport.JsonRpcMessage) error {
	if message.Response != nil {
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
		switch message.Method {
		case mux.RpcRequestMethodCallTool:
			c.handleToolCall(message.Request)
		default:
			c.logger.Error("received message with unexpected method", types.LogArg{
				"method":  message.Method,
				"request": message.Request,
			})
		}
	}
	return nil
}

func (c *ProxyMuxClient) handleProxyRegisterResponse(response *jsonrpc.JsonRpcResponse) error {
	registerResponse, err := mux.ParseJsonRpcResponseProxyRegister(response)
	if err != nil {
		c.logger.Error("error in handleProxyRegisterResponse", types.LogArg{
			"error": err,
		})
		return err
	}

	c.events.EventMuxProxyRegistered(registerResponse)

	return nil
}

func (c *ProxyMuxClient) handleToolCall(request *jsonrpc.JsonRpcRequest) error {
	toolCall, err := mux.ParseJsonRpcRequestToolsCallParams(request)
	if err != nil {
		c.logger.Error("error in handleToolCall", types.LogArg{
			"error": err,
		})
		return err
	}

	c.logger.Info("received tool call", types.LogArg{
		"name":     toolCall.Name,
		"args":     toolCall.Args,
		"mcpReqId": toolCall.McpReqId,
	})

	c.events.EventMuxToolCall(toolCall.Name, toolCall.Args, toolCall.McpReqId)

	return nil
}
