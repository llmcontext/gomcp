package hubmcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

// processing a valid request
func (s *MCPServer) processRequest(ctx context.Context, request *jsonrpc.JsonRpcRequest) error {
	s.logger.Debug("JsonRpcRequest", types.LogArg{
		"request": request,
	})
	switch request.Method {
	case mcp.RpcRequestMethodInitialize:
		{
			parsed, err := mcp.ParseJsonRpcRequestInitialize(request)
			if err != nil {
				s.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
			}
			s.events.EventMcpRequestInitialize(parsed, request.Id)
		}
	case mcp.RpcNotificationMethodInitialized:
		s.events.EventMcpNotificationInitialized()
	case "tools/list":
		return s.handleToolsList(request)
	case mcp.RpcRequestMethodToolsCall:
		return s.handleToolsCall(ctx, request)
	case "resources/list":
		return s.handleResourcesList(request)
	case "prompts/list":
		return s.handlePromptsList(request)
	case "prompts/get":
		return s.handlePromptsGet(request)
	case "ping":
		result := json.RawMessage(`{}`)
		response := &jsonrpc.JsonRpcResponse{
			Id:     request.Id,
			Result: &result,
		}
		return s.sendResponse(response)
	default:
		response := &jsonrpc.JsonRpcResponse{
			Id: request.Id,
			Error: &jsonrpc.JsonRpcError{
				Code:    jsonrpc.RpcMethodNotFound,
				Message: fmt.Sprintf("unknown method: %s", request.Method),
			},
		}
		return s.sendResponse(response)
	}
	return nil
}
