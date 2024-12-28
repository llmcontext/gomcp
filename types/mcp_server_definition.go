package types

type ToolsDefinition interface {
	AddTool(toolName string, description string, toolHandler interface{}) error
}

type McpServerDefinition interface {
	SetDebugLevel(debugLevel string, debugFile string)
	WithTools(configuration interface{}, toolsInitFunction interface{}) ToolsDefinition
}
