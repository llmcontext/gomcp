package tools

type ToolProvider interface {
	AddTool(toolName string, description string, toolHandler interface{}) error
}

type ToolRegistry interface {
	DeclareToolProvider(toolName string, toolInitFunction interface{}) (ToolProvider, error)
}
