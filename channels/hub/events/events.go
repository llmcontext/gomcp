package events

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
)

type Events interface {
	// receive "initialize" request
	EventMcpRequestInitialize(params *mcp.JsonRpcRequestInitializeParams, reqId *jsonrpc.JsonRpcRequestId)

	// receive "initialized" notification
	EventMcpNotificationInitialized()

	// receive "tools/list" request
	EventNewProxyTools()
	// receive "tools/call" request
	EventProxyToolCall(proxyId string, toolName string, toolArgs map[string]interface{}, mcpReqId string)
	// receive "error" notification
	EventMcpError(code int, message string, data *json.RawMessage, id *jsonrpc.JsonRpcRequestId)
}
