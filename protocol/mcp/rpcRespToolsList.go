package mcp

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

type JsonRpcResponseToolsListResult struct {
	Tools      []ToolDescription `json:"tools"`
	NextCursor *string           `json:"nextCursor,omitempty"`
}

type ToolDescription struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

func ParseJsonRpcResponseToolsList(response *jsonrpc.JsonRpcResponse) (*JsonRpcResponseToolsListResult, error) {
	resp := JsonRpcResponseToolsListResult{}

	// parse params
	result, err := protocol.CheckIsObject(response.Result, "result")
	if err != nil {
		return nil, err
	}

	// read tools
	tools, err := protocol.GetArrayField(result, "tools")
	if err != nil {
		return nil, err
	}

	for _, item := range tools {
		tool, err := protocol.CheckIsObject(item, "tool")
		if err != nil {
			return nil, err
		}
		name, err := protocol.GetStringField(tool, "name")
		if err != nil {
			return nil, err
		}

		description, err := protocol.GetStringField(tool, "description")
		if err != nil {
			return nil, err
		}

		inputSchema, err := protocol.GetObjectField(tool, "inputSchema")
		if err != nil {
			return nil, err
		}

		resp.Tools = append(resp.Tools, ToolDescription{
			Name:        name,
			Description: description,
			InputSchema: inputSchema,
		})
	}

	// read next cursor
	nextCursor := protocol.GetOptionalStringField(result, "nextCursor")
	resp.NextCursor = nextCursor

	return &resp, nil
}
