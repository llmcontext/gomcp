package mcpClient

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/jsonrpc/mcp"
)

func (c *MCPProxyClient) handleInitializeResponse(response *jsonrpc.JsonRpcResponse) {
	initializeResponse, err := mcp.ParseJsonRpcResponseInitialize(response)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error in handleInitializeResponse: %+v\n", err))
		return
	}
	c.logger.Info(fmt.Sprintf("init response: %s, %s\n",
		initializeResponse.ServerInfo.Name,
		initializeResponse.ServerInfo.Version,
	))

	// we send the "notifications/initialized" notification
	notification := jsonrpc.NewJsonRpcNotification(mcp.RpcNotificationMethodInitialized)
	c.sendJsonRpcRequest(notification)

	// we send the "tools/list" request
	request, err := mkRpcRequestToolsList(c.clientId)
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to create tools list request: %s\n", err))
		return
	}
	c.sendJsonRpcRequest(request)
}
