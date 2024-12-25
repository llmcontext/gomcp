package mcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

const (
	RpcRequestMethodResourcesList = "resources/list"
)

type JsonRpcRequestResourcesListParams struct {
	Cursor *string `json:"cursor,omitempty"`
}

func ParseJsonRpcRequestResourcesList(params *jsonrpc.JsonRpcParams) (*JsonRpcRequestResourcesListParams, error) {
	resp := &JsonRpcRequestResourcesListParams{}

	// check if we have params
	if params != nil {
		if !params.IsNamed() {
			return nil, fmt.Errorf("invalid call parameters, not an object")
		}
		cursor := protocol.GetOptionalStringField(params.NamedParams, "cursor")
		resp.Cursor = cursor
	}

	return resp, nil
}
