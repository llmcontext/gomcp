package mux

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

const (
	RpcRequestMethodToolsRegister = "tools/register"
)

type JsonRpcRequestToolsRegisterParams struct {
	Tools []ToolDescription `json:"tools"`
}

type ToolDescription struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

func ParseJsonRpcRequestToolsRegisterParams(request *jsonrpc.JsonRpcRequest) (*JsonRpcRequestToolsRegisterParams, error) {
	// parse params
	if request.Params == nil {
		return nil, fmt.Errorf("missing params")
	}
	if !request.Params.IsNamed() {
		return nil, fmt.Errorf("params must be an object")
	}
	namedParams := request.Params.NamedParams

	req := JsonRpcRequestToolsRegisterParams{
		Tools: []ToolDescription{},
	}

	// read tools
	tools, err := protocol.GetArrayField(namedParams, "tools")
	if err != nil {
		return nil, fmt.Errorf("missing tools")
	}
	for _, tool := range tools {
		toolMap, ok := tool.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("tool must be an object")
		}
		toolDescription := ToolDescription{}
		toolDescription.Name, err = protocol.GetStringField(toolMap, "name")
		if err != nil {
			return nil, fmt.Errorf("tool.name must be a string")
		}
		toolDescription.Description, err = protocol.GetStringField(toolMap, "description")
		if err != nil {
			return nil, fmt.Errorf("tool.description must be a string")
		}
		toolDescription.InputSchema, err = protocol.GetObjectField(toolMap, "inputSchema")
		if err != nil {
			return nil, fmt.Errorf("tool.inputSchema must be an object")
		}
		req.Tools = append(req.Tools, toolDescription)
	}

	return &req, nil
}
