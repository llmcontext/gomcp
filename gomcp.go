package gomcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/mcpserver"
	"github.com/llmcontext/gomcp/providers/sdk"
	"github.com/llmcontext/gomcp/types"
)

func NewMcpServerDefinition(serverName string, serverVersion string) types.McpSdkServerDefinition {
	return sdk.NewMcpServerDefinition(serverName, serverVersion)
}

func NewModelContextProtocolServer(definition types.McpSdkServerDefinition) (types.ModelContextProtocolServer, error) {
	mcp, err := mcpserver.NewMcpSdkServer(definition, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create model context protocol server: %v", err)
	}
	return mcp, nil
}
