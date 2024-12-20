package hubmcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
)

func (s *MCPServer) handleToolsCall(ctx context.Context, request *jsonrpc.JsonRpcRequest) error {
	// let's unmarshal the params
	reqParams, err := mcp.ParseJsonRpcRequestToolsCallParams(request.Params)
	if err != nil {
		s.SendError(jsonrpc.RpcInvalidParams, fmt.Sprintf("invalid call parameters: %v", err), request.Id)
		return nil
	}
	toolName := reqParams.Name
	toolArgs := reqParams.Arguments

	// let's check if the tool is a proxy
	isProxy, proxyId, err := s.toolsRegistry.IsProxyTool(ctx, toolName)
	if err != nil {
		s.SendError(jsonrpc.RpcInternalError, fmt.Sprintf("tool call failed: %v", err), request.Id)
		return nil
	}

	// handle proxy tools
	if isProxy {
		s.events.EventProxyToolCall(proxyId, toolName, toolArgs, jsonrpc.RequestIdToString(request.Id))
		return nil
	}

	// let's call the tool
	response, err := s.toolsRegistry.CallTool(ctx, toolName, toolArgs)
	if err != nil {
		s.SendError(jsonrpc.RpcInternalError, fmt.Sprintf("tool call failed: %v", err), request.Id)
		return nil
	}

	// marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		s.SendError(jsonrpc.RpcInternalError, "failed to marshal response", request.Id)
	}
	jsonResponse := json.RawMessage(responseBytes)

	// we send the response
	s.sendResponse(&jsonrpc.JsonRpcResponse{
		Id:     request.Id,
		Result: &jsonResponse,
	})

	return nil
}
