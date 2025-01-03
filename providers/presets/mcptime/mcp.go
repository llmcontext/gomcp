package mcptime

import (
	"context"

	"github.com/llmcontext/gomcp/types"
)

type ToolConfiguration struct {
}

type ToolContext struct {
}

// initializes the McpToolProcessor function.
func ToolInit(ctx context.Context, config *ToolConfiguration) (*ToolContext, error) {
	// we need to initialize the Notion client
	return &ToolContext{}, nil
}

func RegisterTools(mcpServerDefinition types.McpSdkServerDefinition) error {
	mcpToolsDefinition := mcpServerDefinition.WithTools(&ToolConfiguration{}, ToolInit)
	err := mcpToolsDefinition.AddTool("get_current_time", "Get the current time in specified time zone", GetCurrentTime)
	if err != nil {
		return err
	}
	err = mcpToolsDefinition.AddTool("convert_time", "Convert time between two time zones", ConvertTime)
	if err != nil {
		return err
	}

	err = mcpToolsDefinition.AddTool("get_local_timezone", "Get the local timezone", GetLocalTimezone)
	if err != nil {
		return err
	}

	return nil
}
