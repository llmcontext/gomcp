package proxy

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/jsonrpc/messages"
)

func (c *MCPProxyClient) handleToolsListResponse(response *jsonrpc.JsonRpcResponse) {
	toolsListResponse, err := messages.ParseJsonRpcResponseToolsList(response)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error in handleToolsListResponse: %+v\n", err))
		return
	}

	// display all the tool names and descriptions
	for _, tool := range toolsListResponse.Tools {
		c.logger.Info(fmt.Sprintf("tool: %s - %s\n", tool.Name, tool.Description))
	}
}
