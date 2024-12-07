package server

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
)

type promptsListResponse struct {
	Prompts []promptDescription `json:"prompts"`
}

type promptDescription struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Arguments   []argumentDescription `json:"arguments"`
}

type argumentDescription struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

func (s *MCPServer) handlePromptsList(request *jsonrpc.JsonRpcRequest) error {

	var response = promptsListResponse{
		Prompts: make([]promptDescription, 0),
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
