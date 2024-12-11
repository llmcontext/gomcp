package proxy

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
)

func handleInitializeResponse(response *jsonrpc.JsonRpcResponse) {
	fmt.Printf("[proxy] handleInitializeResponse: %+v\n", response)
}
