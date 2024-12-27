package mcp_server_time

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

func RegisterTools(toolRegistry types.ToolRegistry) error {
	toolProvider, err := toolRegistry.DeclareToolProvider("gomcp_server_time", ToolInit)
	if err != nil {
		return err
	}
	err = toolProvider.AddTool("get_current_time", "Get the current time in specified time zone", GetCurrentTime)
	if err != nil {
		return err
	}

	err = toolProvider.AddTool("convert_time", "Convert time between two time zones", ConvertTime)
	if err != nil {
		return err
	}

	err = toolProvider.AddTool("get_local_timezone", "Get the local timezone", GetLocalTimezone)
	if err != nil {
		return err
	}

	return nil
}
