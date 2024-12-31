package presets

import (
	"github.com/llmcontext/gomcp/providers/presets/mcp_server_time"
	"github.com/llmcontext/gomcp/types"
)

func RegisterPresetServers(mcpServerDefinition types.McpSdkServerDefinition, logger types.Logger) error {
	presetToolsNames := []string{
		"gomcp_server_time",
	}

	// TODO: add mechanism to disable some preset tools
	for _, toolName := range presetToolsNames {
		switch toolName {
		case "gomcp_server_time":
			err := mcp_server_time.RegisterTools(mcpServerDefinition)
			if err != nil {
				logger.Error("failed to register tools: %v", types.LogArg{
					"error": err,
				})
			}
		}
	}
	return nil
}
