package mcp

import (
	"github.com/invopop/jsonschema"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/types"
)

// abstract server interface

type PromptArgumentSchema struct {
	Name        string // the name of the argument
	Description string // the description of the argument
	Required    *bool  // true if the argument is required
}

type McpPrompt struct {
	Name        string                 // the name of the prompt
	Description string                 // the description of the prompt
	Arguments   []PromptArgumentSchema // the arguments of the prompt
}

type McpPromptLifecycle struct {
	Init    func()
	Process func(params map[string]string, result types.PromptGetResult, errChan chan *jsonrpc.JsonRpcError)
	End     func()
}

type McpTool struct {
	Name        string             // the name of the tool
	Description string             // the description of the tool
	InputSchema *jsonschema.Schema // A JSON Schema object defining the expected parameters for the tool, top object must be an object.
}

type McpToolLifecycle struct {
	Init    func()
	Process func(params *jsonrpc.JsonRpcParams, result types.ToolCallResult, errChan chan *jsonrpc.JsonRpcError)
	End     func()
}

type McpServer interface {
	AddTool(tool *McpTool, handler McpToolLifecycle) error
	AddPrompt(prompt *McpPrompt, handler McpPromptLifecycle) error
}

type McpServerLifecycle struct {
	Init func()
	End  func()
}

type McpServerRegistry interface {
	RegisterServer(serverName string, serverVersion string, handlers McpServerLifecycle) (McpServer, error)
}
