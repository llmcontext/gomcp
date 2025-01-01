package mcpclient

import (
	"github.com/llmcontext/gomcp/modelcontextprotocol"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type McpClient struct {
	clientName    string
	clientVersion string
	logger        types.Logger
	notifications modelcontextprotocol.McpClientNotifications
	// the transport
	doStopClient     bool
	jsonRpcTransport *transport.JsonRpcTransport
}

func NewMcpClient(
	clientName string,
	clientVersion string,
	notifications modelcontextprotocol.McpClientNotifications,
	logger types.Logger,
) *McpClient {
	return &McpClient{
		clientName:    clientName,
		clientVersion: clientVersion,
		logger:        logger,
		notifications: notifications,
		doStopClient:  false,
	}
}
