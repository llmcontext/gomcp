package sdk

import (
	"fmt"

	"github.com/llmcontext/gomcp/tools"
)

func (s *SdkServerDefinition) CheckConfiguration(toolsRegistry *tools.ToolsRegistry) error {
	// check that the server name and version are not empty
	if s.ServerName() == "" || s.ServerVersion() == "" {
		return fmt.Errorf("invalid MCP server definition: server name or version is empty")
	}

	// declare the tool provider
	// get the type of the configuration
	toolProvider, err := tools.DeclareToolProvider(s.ServerName(), s.toolsInitFunction, s.configuration)
	if err != nil {
		return err
	}
	// we keep track of the tool providers added
	toolsRegistry.RegisterToolProvider(toolProvider)

	// we add all the tools to the tools registry
	for _, tool := range s.toolsDefinition.tools {
		toolProvider.AddTool(tool.toolName, tool.toolDescription, tool.toolHandler)
	}

	return nil
}
