package proxy

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/jsonrpc/messages"
)

func (c *MCPProxyClient) handleIncomingMessage(message jsonrpc.JsonRpcRawMessage, nature jsonrpc.MessageNature) {
	fmt.Printf("[proxy] received message (%d): %+v\n", nature, message)

	if nature == jsonrpc.MessageNatureRequest {
	} else if nature == jsonrpc.MessageNatureResponse {
		response, responseId, err := jsonrpc.ParseJsonRpcResponse(message)
		if err != nil {
			fmt.Printf("[proxy] error parsing response: %+v\n", err)
			return
		}

		if response.Error != nil {
			fmt.Printf("[proxy] error in response: %+v\n", response.Error)
			return
		}

		responseIdString := jsonrpc.RequestIdToString(responseId)

		// we get the pending request
		pendingRequest := c.getPendingRequest(responseId)

		// we check if the pending request is not nil
		if pendingRequest != nil {
			switch pendingRequest.method {
			case messages.RpcRequestMethodInitialize:
				handleInitializeResponse(response)
			}
		} else {
			fmt.Printf("[proxy] no pending request found for response id: %s\n", responseIdString)
		}

		fmt.Printf("[proxy] response: %+v\n", response)
	}
}
