package events

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
)

type Events interface {
	EventNewProxyTools()
	EventProxyToolCall(proxyId string, toolName string, toolArgs map[string]interface{}, mcpReqId string)
	EventMcpError(code int, message string, data *json.RawMessage, id *jsonrpc.JsonRpcRequestId)
}
