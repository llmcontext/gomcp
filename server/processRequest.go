package server

import (
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/logger"
)

// processing a valid request
func (s *MCPServer) processRequest(request *jsonrpc.JsonRpcRequest) error {
	logger.Debug("JsonRpcRequest", logger.Arg{
		"request": request,
	})
	switch request.Method {
	case "initialize":
		return s.handleInitialize(request)
	case "notifications/initialized":
		s.isClientInitialized = true
		// that's a notification, no response is needed
		return nil
	case "tools/list":
		return s.handleToolsList(request)
	case "tools/call":
		return s.handleToolsCall(request)
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
}

func (s *MCPServer) sendResponse(response *jsonrpc.JsonRpcResponse) error {
	logger.Debug("JsonRpcResponse", logger.Arg{
		"response": response,
	})
	jsonResponse, err := jsonrpc.MarshalJsonRpcResponse(response)
	if err != nil {
		return err
	}
	s.transport.Send(jsonResponse)
	return nil
}
