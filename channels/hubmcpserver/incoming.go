package hubmcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (s *MCPServer) handleIncomingMessage(ctx context.Context, message transport.JsonRpcMessage) error {
	if message.Response != nil {
		response := message.Response
		if response.Error != nil {
			s.logger.Error("error in response", types.LogArg{
				"response":      fmt.Sprintf("%+v", response),
				"error_message": response.Error.Message,
				"error_code":    response.Error.Code,
				"error_data":    response.Error.Data,
			})
			return nil
		}
		switch message.Method {
		default:
			s.logger.Error("received message with unexpected method", types.LogArg{
				"method": message.Method,
			})
		}
	} else if message.Request != nil {
		request := message.Request
		switch message.Method {
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
			s.SendJsonRpcResponse(result, request.Id)
		default:
			s.SendError(jsonrpc.RpcMethodNotFound, fmt.Sprintf("unknown method: %s", request.Method), request.Id)
		}
	} else {
		s.logger.Error("received message with unexpected nature", types.LogArg{
			"message": message,
		})
	}

	return nil
}
