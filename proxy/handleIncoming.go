package proxy

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
)

func (c *MCPProxyClient) handleIncomingMessage(message jsonrpc.JsonRpcRawMessage, nature jsonrpc.MessageNature) {
	fmt.Printf("@@ [proxy] received message (%d): %+v\n", nature, message)
}
