package server

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
)

func (s *MCPServer) handleToolsList(request *jsonrpc.JsonRpcRequest) error {
	// we query the tools registry
	tools := s.toolsRegistry.GetListOfTools()

	var response = mcp.JsonRpcResponseToolsListResult{
		Tools: make([]mcp.ToolDescription, 0, len(tools)),
	}

	// we build the response
	for _, tool := range tools {
		// schemaBytes, _ := json.Marshal(tool.InputSchema)
		response.Tools = append(response.Tools, mcp.ToolDescription{
			Name:        tool.ToolName,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}

	// marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInternalError,
			Message: "failed to marshal response",
		}, request.Id)
		return nil
	}
	jsonResponse := json.RawMessage(responseBytes)

	// we send the response
	s.sendResponse(&jsonrpc.JsonRpcResponse{
		Id:     request.Id,
		Result: &jsonResponse,
	})

	return nil
}
