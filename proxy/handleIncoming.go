package proxy

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/jsonrpc/mcp"
)

func (c *MCPProxyClient) handleIncomingMessage(message jsonrpc.JsonRpcRawMessage, nature jsonrpc.MessageNature) {
	if nature == jsonrpc.MessageNatureRequest {
		c.logger.Error(fmt.Sprintf("received JsonRpcRequest: %+v\n", message))
	} else if nature == jsonrpc.MessageNatureResponse {
		response, responseId, err := jsonrpc.ParseJsonRpcResponse(message)
		if err != nil {
			c.logger.Error(fmt.Sprintf("error parsing response: %+v\n", err))
			return
		}

		if response.Error != nil {
			c.logger.Error(fmt.Sprintf("error in response: %+v\n", response.Error))
			return
		}

		responseIdString := jsonrpc.RequestIdToString(responseId)

		// we get the pending request
		pendingRequest := c.getPendingRequest(responseId)

		// we check if the pending request is not nil
		if pendingRequest != nil {
			switch pendingRequest.method {
			case mcp.RpcRequestMethodInitialize:
				c.handleInitializeResponse(response)
			case mcp.RpcRequestMethodToolsList:
				c.handleToolsListResponse(response)
			}
		} else {
			c.logger.Debug(fmt.Sprintf("[proxy] no pending request found for response id: %s\n", responseIdString))
		}
	} else {
		c.logger.Error(fmt.Sprintf("received message with unexpected nature (%d): %+v\n", nature, message))
	}
}
