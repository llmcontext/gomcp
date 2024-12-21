package mux

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

type JsonRpcResponseToolsCallResult struct {
	Content []interface{} `json:"content"`
	IsError *bool         `json:"isError,omitempty"`
}

func ParseJsonRpcResponseToolsCall(response *jsonrpc.JsonRpcResponse) (*JsonRpcResponseToolsCallResult, error) {
	if response.Result == nil {
		return nil, fmt.Errorf("missing result")
	}
	// parse params
	result, err := protocol.CheckIsObject(response.Result, "result")
	if err != nil {
		return nil, err
	}

	// read content
	content, err := protocol.CheckIsArray(result, "content")
	if err != nil {
		return nil, err
	}

	// read isError
	isError := protocol.GetOptionalBoolField(result, "isError")

	return &JsonRpcResponseToolsCallResult{
		Content: content,
		IsError: isError,
	}, nil
}
