package events

import (
	"context"
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/protocol/mux"
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

	// // receive "tools/call" request
	// EventProxyToolCall(proxyId string, toolName string, toolArgs map[string]interface{}, mcpReqId string)

	// receive "resources/list" request
	EventMcpRequestResourcesList(params *mcp.JsonRpcRequestResourcesListParams, reqId *jsonrpc.JsonRpcRequestId)

	// receive "prompts/list" request
	EventMcpRequestPromptsList(params *mcp.JsonRpcRequestPromptsListParams, reqId *jsonrpc.JsonRpcRequestId)

	// receive "prompts/get" request
	EventMcpRequestPromptsGet(params *mcp.JsonRpcRequestPromptsGetParams, reqId *jsonrpc.JsonRpcRequestId)

	// receive "error" notification
	EventMcpError(code int, message string, data *json.RawMessage, id *jsonrpc.JsonRpcRequestId)

	// EventMuxProxyRegister
	EventMuxRequestProxyRegister(proxyId string, params *mux.JsonRpcRequestProxyRegisterParams, reqId *jsonrpc.JsonRpcRequestId)

	// EventMuxRequestToolsRegister
	EventMuxRequestToolsRegister(proxyId string, params *mux.JsonRpcRequestToolsRegisterParams, reqId *jsonrpc.JsonRpcRequestId)

	// EventMuxResponseToolCall
	EventMuxResponseToolCall(toolsCallResult *mux.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId)

	// EventMuxResponseToolCallError
	EventMuxResponseToolCallError(error *jsonrpc.JsonRpcError, reqId *jsonrpc.JsonRpcRequestId)
}
