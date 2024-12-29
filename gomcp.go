package gomcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/channels/hub"
	"github.com/llmcontext/gomcp/sdk"
	"github.com/llmcontext/gomcp/types"
)

func NewMcpServerDefinition(serverName string, serverVersion string) types.McpServerDefinition {
	return sdk.NewMcpServerDefinition(serverName, serverVersion)
}

func NewModelContextProtocolServer(definition types.McpServerDefinition) (types.ModelContextProtocol, error) {
	mcp, err := hub.NewModelContextProtocolServer(definition)
	if err != nil {
		return nil, fmt.Errorf("failed to create model context protocol server: %v", err)
	}
	return mcp, nil
}
