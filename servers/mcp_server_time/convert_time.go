package mcp_server_time

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/llmcontext/gomcp/types"
)

type ConvertTimeInput struct {
	SourceTimeZone string `json:"source_timezone" jsonschema_description:"IANA timezone name (e.g. America/New_York, Europe/London)."`
	Time           string `json:"time" jsonschema_description:"Time to convert in 24-hour format (HH:MM)."`
	TargetTimeZone string `json:"target_timezone" jsonschema_description:"IANA timezone name (e.g. America/New_York, Europe/London)."`
}

type TimeInformation struct {
	TimeZone string `json:"timezone"`
	DateTime string `json:"datetime"`
	IsDST    bool   `json:"is_dst"`
}

type ConvertTimeOutput struct {
	Source         TimeInformation `json:"source"`
	Target         TimeInformation `json:"target"`
	TimeDifference string          `json:"time_difference"`
	DebugLines     []string        `json:"debug_lines"`
}

func ConvertTime(ctx context.Context, toolCtx *ToolContext, input *ConvertTimeInput, output types.ToolCallResult) error {
	var result = ConvertTimeOutput{}
	var debugLines = []string{}
	debugLines = append(debugLines, fmt.Sprintf("Source timezone: %s", input.SourceTimeZone))
	debugLines = append(debugLines, fmt.Sprintf("Target timezone: %s", input.TargetTimeZone))
	debugLines = append(debugLines, fmt.Sprintf("Time to convert: %s", input.Time))

	sourceLocation, err := time.LoadLocation(input.SourceTimeZone)
	if err != nil {
		return fmt.Errorf("invalid source timezone: %s", input.SourceTimeZone)
	}
	debugLines = append(debugLines, fmt.Sprintf("Source location: %s", sourceLocation))

	// let's build a date/time string with the current date and the time to convert
	now := time.Now()
	timeStr := fmt.Sprintf("%s %s", now.Format("2006-01-02"), input.Time)
	sourceTime, err := time.ParseInLocation("2006-01-02 15:04", timeStr, sourceLocation)
	if err != nil {
		return fmt.Errorf("invalid source time: %s", input.Time)
	}
	debugLines = append(debugLines, fmt.Sprintf("Source time: %s", sourceTime))

	targetLocation, err := time.LoadLocation(input.TargetTimeZone)
	if err != nil {
		return fmt.Errorf("invalid target timezone: %s", input.TargetTimeZone)
	}
	debugLines = append(debugLines, fmt.Sprintf("Target location: %s", targetLocation))
	targetTime := sourceTime.In(targetLocation)
	debugLines = append(debugLines, fmt.Sprintf("Target time: %s", targetTime))
	result.Source = TimeInformation{
		TimeZone: sourceLocation.String(),
		DateTime: sourceTime.Format(time.RFC3339),
		IsDST:    isDst(sourceLocation),
	}

	result.Target = TimeInformation{
		TimeZone: targetLocation.String(),
		DateTime: targetTime.Format(time.RFC3339),
		IsDST:    isDst(targetLocation),
	}
	result.DebugLines = debugLines

	// calculate the time difference
	_, sourceOffset := sourceTime.Zone()
	_, targetOffset := targetTime.Zone()
	hoursDiff := (targetOffset - sourceOffset) / 3600
	result.TimeDifference = fmt.Sprintf("%+d hours", hoursDiff)

	// serialize in a formmatted json string
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}

	output.AddTextContent(string(jsonResult))
	return nil
}
