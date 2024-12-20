package events

import "github.com/llmcontext/gomcp/jsonrpc"

type Events interface {
	EventNewProxyTools()
	EventProxyToolCall(proxyId string, toolName string, toolArgs map[string]interface{}, id *jsonrpc.JsonRpcRequestId)
}
