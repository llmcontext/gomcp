package mcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

const (
	RpcRequestMethodToolsList = "tools/list"
)

type JsonRpcRequestToolsListParams struct {
	Cursor *string `json:"cursor,omitempty"`
}

func ParseJsonRpcRequestToolsList(request *jsonrpc.JsonRpcRequest) (*JsonRpcRequestToolsListParams, error) {
	params := &JsonRpcRequestToolsListParams{}

	// check if we have params
	if request.Params != nil {
		if !request.Params.IsNamed() {
			return nil, fmt.Errorf("invalid call parameters, not an object")
		}
		cursor := protocol.GetOptionalStringField(request.Params.NamedParams, "cursor")
		params.Cursor = cursor
	}

	return params, nil
}
