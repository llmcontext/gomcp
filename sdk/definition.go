package sdk

import (
	"slices"

	"github.com/llmcontext/gomcp/types"
)

type SdkServerDefinition struct {
	serverName        string
	serverVersion     string
	debugLevel        string
	debugFile         string
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

func (s *SdkServerDefinition) SetDebugLevel(debugLevel string, debugFile string) {
	s.debugLevel = debugLevel
	s.debugFile = debugFile
}

func (s *SdkServerDefinition) DebugLevel() string {
	validLevels := []string{"debug", "info", "warn", "error", "dpanic", "panic", "fatal"}
	if !slices.Contains(validLevels, s.debugLevel) {
		return "info"
	}
	return s.debugLevel
}

func (s *SdkServerDefinition) DebugFile() string {
	if s.debugFile == "" {
		return ""
	}
	return s.debugFile
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

func (s *SdkServerDefinition) serverInitFunction() {

}

func (s *SdkServerDefinition) serverEndFunction() {

}
