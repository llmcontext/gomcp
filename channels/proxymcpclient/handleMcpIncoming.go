package proxymcpclient

import (
	"fmt"

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
				"response":      fmt.Sprintf("%+v", response),
				"error":         response.Error,
				"error_message": response.Error.Message,
				"error_code":    response.Error.Code,
				"error_data":    response.Error.Data,
			})
			return
		}
		// we check if the pending request is not nil
		switch message.Method {
		case mcp.RpcRequestMethodInitialize:
			c.handleMcpInitializeResponse(response)
		case mcp.RpcRequestMethodToolsList:
			c.handleMcpToolsListResponse(response)
		case mcp.RpcRequestMethodToolsCall:
			c.handleMcpToolsCallResponse(response, message.ExtraParam)
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
