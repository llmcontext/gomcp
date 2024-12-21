package proxymcpclient

import (
	"fmt"

	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMcpClient) handleIncomingMessage(message transport.JsonRpcMessage) error {
	if message.Response != nil {
		response := message.Response
		if response.Error != nil {
			c.logger.Error("error in response", types.LogArg{
				"response":      fmt.Sprintf("%+v", response),
				"error_message": response.Error.Message,
				"error_code":    response.Error.Code,
				"error_data":    response.Error.Data,
			})
			return nil
		}
		switch message.Method {
		case mcp.RpcRequestMethodInitialize:
			{
				initializeResponse, err := mcp.ParseJsonRpcResponseInitialize(response)
				if err != nil {
					c.logger.Error("error in handleMcpInitializeResponse", types.LogArg{
						"error": err,
					})
					return nil
				}
				c.logger.Info("init response", types.LogArg{
					"name":    initializeResponse.ServerInfo.Name,
					"version": initializeResponse.ServerInfo.Version,
				})
				c.events.EventMcpResponseInitialize(initializeResponse)
			}
		case mcp.RpcRequestMethodToolsList:
			c.handleMcpToolsListResponse(response)
		case mcp.RpcRequestMethodToolsCall:
			c.handleMcpToolsCallResponse(response, message.ExtraParam)
		default:
			c.logger.Error("received message with unexpected method", types.LogArg{
				"method": message.Method,
			})
		}
	} else if message.Request != nil {
		request := message.Request
		switch message.Method {
		default:
			c.logger.Error("received message with unexpected method", types.LogArg{
				"method":  message.Method,
				"request": request,
			})
		}
	} else {
		c.logger.Error("received message with unexpected nature", types.LogArg{
			"message": message,
		})
	}

	return nil
}
