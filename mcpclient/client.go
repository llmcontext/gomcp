package mcpclient

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type McpClient struct {
	clientName    string
	clientVersion string
	logger        types.Logger
	notifications McpClientNotifications
	// the transport
	jsonRpcTransport *transport.JsonRpcTransport
}

type McpClientNotifications interface {
	OnServerInformation(serverName string, serverVersion string)
	OnToolsList(result *mcp.JsonRpcResponseToolsListResult, rpcError *jsonrpc.JsonRpcError)
	OnToolCallResponse(result *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, rpcError *jsonrpc.JsonRpcError)
}

func NewMcpClient(clientName string, clientVersion string, logger types.Logger, notifications McpClientNotifications) *McpClient {
	return &McpClient{
		clientName:    clientName,
		clientVersion: clientVersion,
		logger:        logger,
		notifications: notifications,
	}
}
