package sdk

import (
	"reflect"
	"slices"

	"github.com/invopop/jsonschema"
	"github.com/llmcontext/gomcp/types"
)

type SdkServerDefinition struct {
	serverName            string
	serverVersion         string
	debugLevel            string
	debugFile             string
	toolConfigurationData interface{}
	toolsInitFunction     interface{}
	toolDefinitions       []*SdkToolDefinition

	// enhanced data
	contextType     reflect.Type
	contextTypeName string
	// the tool context retrieve from the tool init function
	toolContext interface{}
}

type SdkToolDefinition struct {
	ToolName            string
	ToolDescription     string
	toolHandlerFunction interface{}

	// from the server context
	toolContext interface{}

	// enhanced data
	InputSchema   *jsonschema.Schema
	inputTypeName string
}

func NewMcpSdkServerDefinition(serverName string, serverVersion string) *SdkServerDefinition {
	return &SdkServerDefinition{
		serverName:      serverName,
		serverVersion:   serverVersion,
		toolDefinitions: []*SdkToolDefinition{},
	}
}

func NewMcpServerDefinition(serverName string, serverVersion string) types.McpSdkServerDefinition {
	return NewMcpSdkServerDefinition(serverName, serverVersion).AsMcpServerDefinition()
}

func (s *SdkServerDefinition) AsMcpServerDefinition() types.McpSdkServerDefinition {
	return s
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

func (s *SdkServerDefinition) WithTools(toolConfigurationDate interface{}, toolsInitFunction interface{}) types.ToolsDefinition {
	s.toolConfigurationData = toolConfigurationDate
	s.toolsInitFunction = toolsInitFunction
	return s
}

func (s *SdkServerDefinition) AddTool(toolName string, description string, toolHandler interface{}) error {
	s.toolDefinitions = append(s.toolDefinitions, &SdkToolDefinition{
		ToolName:            toolName,
		ToolDescription:     description,
		toolHandlerFunction: toolHandler,
	})
	return nil
}

func (s *SdkServerDefinition) GetListOfTools() []*SdkToolDefinition {
	return s.toolDefinitions
}

func (s *SdkServerDefinition) GetTool(toolName string) *SdkToolDefinition {
	for _, tool := range s.toolDefinitions {
		if tool.ToolName == toolName {
			return tool
		}
	}
	return nil
}
