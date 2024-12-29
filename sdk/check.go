package sdk

import (
	"fmt"

	"github.com/llmcontext/gomcp/mcp"
	"github.com/llmcontext/gomcp/tools"
)

func (s *SdkServerDefinition) SetupMcpServer(
	toolsRegistry *tools.ToolsRegistry,
	mcpServerRegistry mcp.McpServerRegistry) error {
	// check that the server name and version are not empty
	if s.ServerName() == "" || s.ServerVersion() == "" {
		return fmt.Errorf("invalid MCP server definition: server name or version is empty")
	}

	// server handlers
	serverHandlers := mcp.McpServerLifecycle{
		Init: s.serverInitFunction,
		End:  s.serverEndFunction,
	}
	mcpServer, err := mcpServerRegistry.RegisterServer(s.ServerName(), s.ServerVersion(), serverHandlers)
	if err != nil {
		return fmt.Errorf("failed to register MCP server: %v", err)
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
		toolLifecycle := mcp.McpToolLifecycle{
			Init:    nil,
			Process: nil,
			End:     nil,
		}
		toolDefinition := mcp.McpTool{
			Name:        tool.toolName,
			Description: tool.toolDescription,
			InputSchema: nil,
		}
		mcpServer.AddTool(&toolDefinition, toolLifecycle)

		toolProvider.AddTool(tool.toolName, tool.toolDescription, tool.toolHandler)
	}

	return nil
}
