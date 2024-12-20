package mux

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

type JsonRpcResponseToolsCallResult struct {
	Content  []interface{} `json:"content"`
	IsError  *bool         `json:"isError,omitempty"`
	McpReqId string        `json:"mcp_req_id"`
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

	// read mcpReqId
	mcpReqId, err := protocol.GetStringField(result, "mcp_req_id")
	if err != nil {
		return nil, err
	}

	return &JsonRpcResponseToolsCallResult{
		Content:  content,
		IsError:  isError,
		McpReqId: mcpReqId,
	}, nil
}
