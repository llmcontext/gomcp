package mcpclient

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (c *MCPProxyClient) handleMcpInitializeResponse(
	response *jsonrpc.JsonRpcResponse,
	transport *transport.JsonRpcTransport,
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

	// we update the server information
	c.serverInfo.Name = initializeResponse.ServerInfo.Name
	c.serverInfo.Version = initializeResponse.ServerInfo.Version

	// we send the "notifications/initialized" notification
	notification := jsonrpc.NewJsonRpcNotification(mcp.RpcNotificationMethodInitialized)
	transport.SendRequest(notification)

	// we send the "tools/list" request
	params := mcp.JsonRpcRequestToolsListParams{}

	transport.SendRequestWithMethodAndParams(
		mcp.RpcRequestMethodToolsList, params)
}
