package hubmcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
)

func (s *MCPServer) handleToolsCall(ctx context.Context, request *jsonrpc.JsonRpcRequest) error {
	// let's unmarshal the params
	reqParams := request.Params

	if reqParams == nil || reqParams.NamedParams == nil {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidParams,
			Message: "invalid call parameters, not an object",
		}, request.Id)
		return nil
	}

	// let's get the tool name
	// we expect a string
	toolName, ok := reqParams.NamedParams["name"].(string)
	if !ok {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidParams,
			Message: "invalid tool name",
		}, request.Id)
		return nil
	}

	// let's get the tool arguments
	// we expect a map[string]interface{}
	toolArgs, ok := reqParams.NamedParams["arguments"].(map[string]interface{})
	if !ok {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidParams,
			Message: "invalid tool arguments",
		}, request.Id)
		return nil
	}

	// let's call the tool
	response, err := s.toolsRegistry.CallTool(ctx, toolName, toolArgs)
	if err != nil {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInternalError,
			Message: fmt.Sprintf("tool call failed: %v", err),
		}, request.Id)
		return nil
	}

	// marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInternalError,
			Message: "failed to marshal response",
		}, request.Id)
	}
	jsonResponse := json.RawMessage(responseBytes)

	// we send the response
	s.sendResponse(&jsonrpc.JsonRpcResponse{
		Id:     request.Id,
		Result: &jsonResponse,
	})

	return nil
}
