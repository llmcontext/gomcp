package proxymcpclient

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMcpClient) handleMcpToolsListResponse(
	response *jsonrpc.JsonRpcResponse,
) {
	// the MCP server sent its tools list
	toolsListResponse, err := mcp.ParseJsonRpcResponseToolsList(response)
	if err != nil {
		c.logger.Error("error in handleMcpToolsListResponse", types.LogArg{
			"error": err,
		})
		return
	}

	c.events.EventMcpToolsListResponse(toolsListResponse)

}
