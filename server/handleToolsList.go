package server

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
)

type toolsListResponse struct {
	Tools []toolDescription `json:"tools"`
}

type toolDescription struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

func (s *MCPServer) handleToolsList(request *jsonrpc.JsonRpcRequest) error {
	// we query the tools registry
	tools := s.toolsRegistry.GetListOfTools()

	var response = toolsListResponse{
		Tools: make([]toolDescription, 0, len(tools)),
	}

	// we build the response
	for _, tool := range tools {
		// schemaBytes, _ := json.Marshal(tool.InputSchema)
		response.Tools = append(response.Tools, toolDescription{
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
