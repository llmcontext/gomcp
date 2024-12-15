package mcpClient

import (
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (c *MCPProxyClient) handleMcpIncomingMessage(message transport.JsonRpcMessage,
	transport *transport.JsonRpcTransport) {
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

		// we get the pending request
		pendingRequestMethod, reqId := transport.GetPendingRequest(response.Id)
		if reqId == nil {
			c.logger.Debug("[proxy] no pending request found for response id", types.LogArg{
				"response": response,
			})
			return
		}

		// we check if the pending request is not nil
		if pendingRequestMethod != "" {
			switch pendingRequestMethod {
			case mcp.RpcRequestMethodInitialize:
				c.handleMcpInitializeResponse(response, transport)
			case mcp.RpcRequestMethodToolsList:
				c.handleMcpToolsListResponse(response, transport)
			}
		} else {
			c.logger.Debug("[proxy] no pending request found for response id", types.LogArg{
				"id": response.Id,
			})
		}
	} else {
		c.logger.Error("received message with unexpected nature", types.LogArg{
			"message": message,
		})
	}
}
