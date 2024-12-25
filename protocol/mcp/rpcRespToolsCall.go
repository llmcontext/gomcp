package mcp

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

type JsonRpcResponseToolsCallResult struct {
	Content []interface{} `json:"content"`
	IsError *bool         `json:"isError,omitempty"`
}

func ParseJsonRpcResponseToolsCall(response *jsonrpc.JsonRpcResponse) (*JsonRpcResponseToolsCallResult, error) {
	resp := JsonRpcResponseToolsCallResult{
		Content: []interface{}{},
		IsError: nil,
	}

	// parse params
	result, err := protocol.CheckIsObject(response.Result, "result")
	if err != nil {
		return nil, err
	}

	// parse params
	content, err := protocol.CheckIsArray(result["content"], "content")
	if err != nil {
		return nil, err
	}

	resp.Content = content

	isError := protocol.GetOptionalBoolField(result, "isError")
	if isError != nil {
		resp.IsError = isError
	}

	return &resp, nil
}
