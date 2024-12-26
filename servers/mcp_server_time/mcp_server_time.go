package mcp_server_time

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

type GetCurrentTimeInput struct {
	TimeZone string `json:"timezone"`
}

type GetCurrentTimeOutput struct {
	TimeZone    string `json:"timezone"`
	DateTime    string `json:"datetime"`
	SupportsDST bool   `json:"is_dst"`
}

func GetCurrentTime(ctx context.Context, toolCtx *ToolContext, input *GetCurrentTimeInput, output types.ToolCallResult) error {
	// Load location for the given timezone
	location, err := time.LoadLocation(input.TimeZone)
	if err != nil {
		return fmt.Errorf("invalid timezone: %s", input.TimeZone)
	}

	// Get current time in the specified timezone
	now := time.Now().In(location)

	// Check if the timezone has DST
	_, winterOffset := now.Zone()
	summerTime := now.AddDate(0, 6, 0) // Check 6 months later
	_, summerOffset := summerTime.Zone()
	hasDST := winterOffset != summerOffset

	result := GetCurrentTimeOutput{
		TimeZone:    location.String(),
		DateTime:    now.Format(time.RFC3339),
		SupportsDST: hasDST,
	}

	// serialize the result
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return err
	}

	output.AddTextContent(string(jsonResult))
	return nil
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
	return nil
}
