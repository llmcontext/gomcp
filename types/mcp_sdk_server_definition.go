package types

import "github.com/llmcontext/gomcp/pkg/prompts"

type ToolsDefinition interface {
	AddTool(toolName string, description string, toolHandler interface{}) error
}

type McpSdkServerDefinition interface {
	SetDebugLevel(debugLevel string, debugFile string)
	WithTools(configuration interface{}, toolsInitFunction interface{}) ToolsDefinition
	AddTemplateYamlFile(templateYamlFilePath string) ([]*prompts.DuplicatedPrompt, error)
}
