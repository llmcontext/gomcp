package server

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
)

type resourcesListResponse struct {
	Resources []resourceDescription `json:"resources"`
}

type resourceDescription struct {
	Uri         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description"`
	MimeType    string `json:"mimeType"`
}

func (s *MCPServer) handleResourcesList(request *jsonrpc.JsonRpcRequest) error {

	var response = resourcesListResponse{
		Resources: make([]resourceDescription, 0),
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
