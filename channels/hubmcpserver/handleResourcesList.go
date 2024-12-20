package hubmcpserver

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
)

func (s *MCPServer) handleResourcesList(request *jsonrpc.JsonRpcRequest) error {

	var response = mcp.JsonRpcResponseResourcesListResult{
		Resources: make([]mcp.ResourceDescription, 0),
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
