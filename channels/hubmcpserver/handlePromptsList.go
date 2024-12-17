package hubmcpserver

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
)

func (s *MCPServer) handlePromptsList(request *jsonrpc.JsonRpcRequest) error {

	var response = mcp.JsonRpcResponsePromptsListResult{
		Prompts: make([]mcp.PromptDescription, 0),
	}

	prompts := s.promptsRegistry.GetListOfPrompts()
	for _, prompt := range prompts {
		arguments := make([]mcp.PromptArgumentDescription, 0, len(prompt.Arguments))
		for _, argument := range prompt.Arguments {
			arguments = append(arguments, mcp.PromptArgumentDescription{
				Name:        argument.Name,
				Description: argument.Description,
				Required:    argument.Required,
			})
		}
		response.Prompts = append(response.Prompts, mcp.PromptDescription{
			Name:        prompt.Name,
			Description: prompt.Description,
			Arguments:   arguments,
		})
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
