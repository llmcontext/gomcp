package proxymcpclient

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMcpClient) handleMcpInitializeResponse(
	response *jsonrpc.JsonRpcResponse,
) {
	initializeResponse, err := mcp.ParseJsonRpcResponseInitialize(response)
	if err != nil {
		c.logger.Error("error in handleMcpInitializeResponse", types.LogArg{
			"error": err,
		})
		return
	}
	c.logger.Info("init response", types.LogArg{
		"name":    initializeResponse.ServerInfo.Name,
		"version": initializeResponse.ServerInfo.Version,
	})
	c.events.EventMcpInitializeResponse(initializeResponse)

}
