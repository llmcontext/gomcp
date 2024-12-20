package proxymcpclient

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMcpClient) handleMcpToolsCallResponse(response *jsonrpc.JsonRpcResponse, extraParam string) {
	// we got the response from the mcp server

	toolsCallResult, err := mcp.ParseJsonRpcResponseToolsCall(response)
	if err != nil {
		c.logger.Error("error parsing tools call params", types.LogArg{
			"error": err,
		})
		return
	}

	c.logger.Info("tools call result", types.LogArg{
		"content":    toolsCallResult.Content,
		"isError":    toolsCallResult.IsError,
		"extraParam": extraParam,
	})

	// TODO: we need to forward the response to the hubmux server
}
