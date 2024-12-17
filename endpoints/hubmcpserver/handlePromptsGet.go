package hubmcpserver

import (
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
)

type GetPromptResult struct {
	Description string          `json:"description"`
	Messages    []PromptMessage `json:"messages"`
}

type PromptMessage struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

func (s *MCPServer) handlePromptsGet(request *jsonrpc.JsonRpcRequest) error {
	// let's unmarshal the params
	reqParams := request.Params

	if reqParams == nil || reqParams.NamedParams == nil {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidParams,
			Message: "invalid call parameters, not an object",
		}, request.Id)
		return nil
	}

	promptName, ok := reqParams.NamedParams["name"].(string)
	if !ok {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidParams,
			Message: "invalid call parameters, missing 'name'",
		}, request.Id)
		return nil
	}

	var templateArgs = map[string]string{}
	// let's validate the arguments
	if reqParams.NamedParams["arguments"] != nil {
		args, ok := reqParams.NamedParams["arguments"].(map[string]interface{})
		if !ok {
			s.sendError(&jsonrpc.JsonRpcError{
				Code:    jsonrpc.RpcInvalidParams,
				Message: "invalid call parameters, 'arguments' is not an object",
			}, request.Id)
			return nil
		}
		// copy the arguments, as strings
		for key, value := range args {
			templateArgs[key] = fmt.Sprintf("%v", value)
		}
	}

	response, err := s.promptsRegistry.GetPrompt(promptName, templateArgs)
	if err != nil {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidParams,
			Message: fmt.Sprintf("prompt processing error: %s", err),
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
