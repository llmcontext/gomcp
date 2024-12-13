package mcp

import (
	"github.com/llmcontext/gomcp/jsonrpc"
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
	result, err := checkIsObject(response.Result, "result")
	if err != nil {
		return nil, err
	}

	// read tools
	tools, err := getArrayField(result, "tools")
	if err != nil {
		return nil, err
	}

	for _, item := range tools {
		tool, err := checkIsObject(item, "tool")
		if err != nil {
			return nil, err
		}
		name, err := getStringField(tool, "name")
		if err != nil {
			return nil, err
		}

		description, err := getStringField(tool, "description")
		if err != nil {
			return nil, err
		}

		inputSchema, err := getObjectField(tool, "inputSchema")
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
	nextCursor := getOptionalStringField(result, "nextCursor")
	resp.NextCursor = nextCursor

	return &resp, nil
}
