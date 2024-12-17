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

	// display all the tool names and descriptions
	for _, tool := range toolsListResponse.Tools {
		c.logger.Info("tool", types.LogArg{
			"name":        tool.Name,
			"description": tool.Description,
		})
		// we store the tools description
		c.tools = append(c.tools, tool)
	}

	// we can now report the tools list to the mux server
	c.muxClient.SendProxyRegistrationRequest(c.options, c.serverInfo, c.tools)
}
