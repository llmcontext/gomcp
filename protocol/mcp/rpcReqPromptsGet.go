package mcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

const (
	RpcRequestMethodPromptsGet = "prompts/get"
)

type JsonRpcRequestPromptsGetParams struct {
	Name      string            `json:"name"`
	Arguments map[string]string `json:"arguments"`
}

func ParseJsonRpcRequestPromptsGet(params *jsonrpc.JsonRpcParams) (*JsonRpcRequestPromptsGetParams, error) {
	resp := &JsonRpcRequestPromptsGetParams{}

	// check if we have params
	if params == nil {
		return nil, fmt.Errorf("invalid call parameters, no parameters provided")
	}

	if !params.IsNamed() {
		return nil, fmt.Errorf("invalid call parameters, not an object")
	}

	name, err := protocol.GetStringField(params.NamedParams, "name")
	if err != nil {
		return nil, fmt.Errorf("invalid call parameters, name is not a string")
	}
	resp.Name = name

	arguments, err := protocol.GetObjectField(params.NamedParams, "arguments")
	if err != nil {
		return nil, fmt.Errorf("invalid call parameters, arguments is not an object")
	}

	// let's convert the arguments to a map[string]string
	resp.Arguments = make(map[string]string)
	for k, v := range arguments {
		resp.Arguments[k] = v.(string)
	}

	// check if the request is valid
	return resp, nil
}
