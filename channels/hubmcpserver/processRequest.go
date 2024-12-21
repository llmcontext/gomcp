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
	case mcp.RpcRequestMethodToolsList:
		{
			parsed, err := mcp.ParseJsonRpcRequestToolsList(request)
			if err != nil {
				s.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
			}
			s.events.EventMcpRequestToolsList(parsed, request.Id)
		}
	case mcp.RpcRequestMethodToolsCall:
		{
			parsed, err := mcp.ParseJsonRpcRequestToolsCallParams(request.Params)
			if err != nil {
				s.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
				return nil
			}
			s.events.EventMcpRequestToolsCall(ctx, parsed, request.Id)
		}
	case mcp.RpcRequestMethodResourcesList:
		{
			parsed, err := mcp.ParseJsonRpcRequestResourcesList(request.Params)
			if err != nil {
				s.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
			}
			s.events.EventMcpRequestResourcesList(parsed, request.Id)
		}
	case mcp.RpcRequestMethodPromptsList:
		{
			parsed, err := mcp.ParseJsonRpcRequestPromptsList(request.Params)
			if err != nil {
				s.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
			}
			s.events.EventMcpRequestPromptsList(parsed, request.Id)
		}
	case mcp.RpcRequestMethodPromptsGet:
		{
			parsed, err := mcp.ParseJsonRpcRequestPromptsGet(request.Params)
			if err != nil {
				s.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
			}
			s.events.EventMcpRequestPromptsGet(parsed, request.Id)
		}
	case "ping":
		result := json.RawMessage(`{}`)
		response := &jsonrpc.JsonRpcResponse{
			Id:     request.Id,
			Result: &result,
		}
		s.SendResponse(response)
	default:
		s.SendError(jsonrpc.RpcMethodNotFound, fmt.Sprintf("unknown method: %s", request.Method), request.Id)
	}
	return nil
}
