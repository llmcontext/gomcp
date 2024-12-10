package proxy

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
)

func (c *MCPProxyClient) handleIncomingMessage(message jsonrpc.JsonRpcRawMessage, nature jsonrpc.MessageNature) {
	fmt.Printf("@@ [proxy] received message (%d): %+v\n", nature, message)

	if nature == jsonrpc.MessageNatureRequest {
	} else if nature == jsonrpc.MessageNatureResponse {
		response, responseId, err := jsonrpc.ParseJsonRpcResponse(message)
		if err != nil {
			fmt.Printf("@@ [proxy] error parsing response: %+v\n", err)
			return
		}

		if response.Error != nil {
			fmt.Printf("@@ [proxy] error in response: %+v\n", response.Error)
			return
		}

		if responseId != nil {
			fmt.Printf("@@ [proxy] response id: %+v\n", responseId)
		}

		fmt.Printf("@@ [proxy] response: %+v\n", response)
	}
}
