package mcp_server_time

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/llmcontext/gomcp/types"
)

type GetLocalTimezoneInput struct {
}

func getIANATimezone() (string, error) {
	// Check if /var/db/timezone/tz/*/zoneinfo is a symlink
	tzPath, err := filepath.EvalSymlinks("/etc/localtime")
	if err != nil {
		return "", err
	}

	// On macOS, the path looks like: /private/var/db/timezone/tz/2024a.1.0/zoneinfo/America/Los_Angeles
	if strings.Contains(tzPath, "/var/db/timezone/tz/") {
		parts := strings.Split(tzPath, "/zoneinfo/")
		if len(parts) == 2 {
			return parts[1], nil
		}
	}

	// For Linux systems, check the traditional path
	const zoneinfo = "/usr/share/zoneinfo/"
	if strings.HasPrefix(tzPath, zoneinfo) {
		return tzPath[len(zoneinfo):], nil
	}

	return "", fmt.Errorf("unknown timezone path: %s", tzPath)
}

func GetLocalTimezone(ctx context.Context, toolCtx *ToolContext, input *GetLocalTimezoneInput, output types.ToolCallResult) error {
	ianaTimezone, err := getIANATimezone()
	if err != nil {
		return err
	}

	output.AddTextContent(ianaTimezone)
	return nil
}
