package events

import (
	"context"
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
	EventMcpRequestToolsList(params *mcp.JsonRpcRequestToolsListParams, reqId *jsonrpc.JsonRpcRequestId)

	// receive "tools/call" request
	EventMcpRequestToolsCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams, reqId *jsonrpc.JsonRpcRequestId)

	// receive "tools/call" request
	EventProxyToolCall(proxyId string, toolName string, toolArgs map[string]interface{}, mcpReqId string)

	// receive "resources/list" request
	EventMcpRequestResourcesList(params *mcp.JsonRpcRequestResourcesListParams, reqId *jsonrpc.JsonRpcRequestId)

	// receive "prompts/list" request
	EventMcpRequestPromptsList(params *mcp.JsonRpcRequestPromptsListParams, reqId *jsonrpc.JsonRpcRequestId)

	// receive "prompts/get" request
	EventMcpRequestPromptsGet(params *mcp.JsonRpcRequestPromptsGetParams, reqId *jsonrpc.JsonRpcRequestId)

	// receive "error" notification
	EventMcpError(code int, message string, data *json.RawMessage, id *jsonrpc.JsonRpcRequestId)
}
