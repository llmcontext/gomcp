package sdk

import "github.com/llmcontext/gomcp/types"

type SdkServerDefinition struct {
	serverName        string
	serverVersion     string
	configuration     interface{}
	toolsInitFunction interface{}
	toolsDefinition   SdkToolsDefinition
}

type SdkToolsDefinition struct {
	tools []SdkToolDefinition
}

type SdkToolDefinition struct {
	toolName        string
	toolDescription string
	toolHandler     interface{}
}

func NewMcpServerDefinition(serverName string, serverVersion string) types.McpServerDefinition {
	return &SdkServerDefinition{
		serverName:    serverName,
		serverVersion: serverVersion,
		toolsDefinition: SdkToolsDefinition{
			tools: []SdkToolDefinition{},
		},
	}
}

func (s *SdkServerDefinition) ServerName() string {
	return s.serverName
}

func (s *SdkServerDefinition) ServerVersion() string {
	return s.serverVersion
}

func (s *SdkServerDefinition) WithTools(configuration interface{}, toolsInitFunction interface{}) types.ToolsDefinition {
	s.configuration = configuration
	s.toolsInitFunction = toolsInitFunction
	return &s.toolsDefinition
}

func (s *SdkToolsDefinition) AddTool(toolName string, description string, toolHandler interface{}) error {
	s.tools = append(s.tools, SdkToolDefinition{
		toolName:        toolName,
		toolDescription: description,
		toolHandler:     toolHandler,
	})
	return nil
}
