package mcp_server_time

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/llmcontext/gomcp/types"
)

type GetCurrentTimeInput struct {
	TimeZone string `json:"timezone" jsonschema_description:"IANA timezone name (e.g. America/New_York, Europe/London)."`
}

type GetCurrentTimeOutput struct {
	TimeZone string `json:"timezone"`
	DateTime string `json:"datetime"`
	IsDST    bool   `json:"is_dst"`
}

func GetCurrentTime(ctx context.Context, toolCtx *ToolContext, input *GetCurrentTimeInput, output types.ToolCallResult) error {
	// Load location for the given timezone
	location, err := time.LoadLocation(input.TimeZone)
	if err != nil {
		return fmt.Errorf("invalid timezone: %s", input.TimeZone)
	}

	// Get current time in UTC first, then convert to desired timezone
	now := time.Now().UTC().In(location)

	// Get current name of timezone (handles DST automatically)
	_, offset := now.Zone()

	result := GetCurrentTimeOutput{
		TimeZone: location.String(),
		DateTime: now.Format(time.RFC3339),
		IsDST:    offset != 0 && isDst(location),
	}

	// serialize the result
	jsonResult, err := json.Marshal(result)
	if err != nil {
		return err
	}

	output.AddTextContent(string(jsonResult))
	return nil
}
