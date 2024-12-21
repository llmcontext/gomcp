package mcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

const (
	RpcRequestMethodPromptsList = "prompts/list"
)

type JsonRpcRequestPromptsListParams struct {
	Cursor *string `json:"cursor,omitempty"`
}

func ParseJsonRpcRequestPromptsList(params *jsonrpc.JsonRpcParams) (*JsonRpcRequestPromptsListParams, error) {
	resp := &JsonRpcRequestPromptsListParams{}

	// check if we have params
	if params != nil {
		if !params.IsNamed() {
			return nil, fmt.Errorf("invalid call parameters, not an object")
		}
		cursor, err := protocol.GetStringField(params.NamedParams, "cursor")
		if err != nil {
			return nil, fmt.Errorf("invalid call parameters, cursor is not a string")
		}
		resp.Cursor = &cursor
	}

	return resp, nil
}
