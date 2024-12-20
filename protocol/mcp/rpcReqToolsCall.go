package mcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

const (
	RpcRequestMethodToolsCall = "tools/call"
)

type JsonRpcRequestToolsCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

func ParseJsonRpcRequestToolsCallParams(params *jsonrpc.JsonRpcParams) (*JsonRpcRequestToolsCallParams, error) {
	if params == nil {
		return nil, fmt.Errorf("invalid call parameters, not an object")
	}
	if !params.IsNamed() {
		return nil, fmt.Errorf("params must be an object")
	}
	namedParams := params.NamedParams

	toolCall := &JsonRpcRequestToolsCallParams{}

	// check if name is present
	name, err := protocol.GetStringField(namedParams, "name")
	if err != nil {
		return nil, fmt.Errorf("missing name")
	}
	toolCall.Name = name

	// check if args is present
	arguments, err := protocol.GetObjectField(namedParams, "arguments")
	if err != nil {
		return nil, fmt.Errorf("missing arguments")
	}
	toolCall.Arguments = arguments

	return toolCall, nil
}
