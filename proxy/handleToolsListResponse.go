package proxy

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/jsonrpc/mcp"
)

func (c *MCPProxyClient) handleToolsListResponse(response *jsonrpc.JsonRpcResponse) {
	toolsListResponse, err := mcp.ParseJsonRpcResponseToolsList(response)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error in handleToolsListResponse: %+v\n", err))
		return
	}

	// display all the tool names and descriptions
	for _, tool := range toolsListResponse.Tools {
		c.logger.Info(fmt.Sprintf("tool: %s - %s\n", tool.Name, tool.Description))
		c.logger.Info(fmt.Sprintf("inputSchema: %+#v\n", tool.InputSchema))
	}
}
