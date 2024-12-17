package proxymcpclient

import (
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMcpClient) handleMcpIncomingMessage(message transport.JsonRpcMessage) {
	if message.Request != nil {
		c.logger.Error("received JsonRpcRequest", types.LogArg{
			"request": message.Request,
		})
	} else if message.Response != nil {
		response := message.Response
		if response.Error != nil {
			c.logger.Error("error in response", types.LogArg{
				"response": response,
			})
			return
		}
		// we check if the pending request is not nil
		switch message.Method {
		case mcp.RpcRequestMethodInitialize:
			c.handleMcpInitializeResponse(response)
		case mcp.RpcRequestMethodToolsList:
			c.handleMcpToolsListResponse(response)
		default:
			c.logger.Error("received message with unexpected method", types.LogArg{
				"method": message.Method,
			})
		}
	} else {
		c.logger.Error("received message with unexpected nature", types.LogArg{
			"message": message,
		})
	}
}
