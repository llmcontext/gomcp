package mcpserver

import (
	"fmt"

	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/modelcontextprotocol"
	"github.com/llmcontext/gomcp/providers"
	"github.com/llmcontext/gomcp/providers/sdk"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type McpServer struct {
	logger        types.Logger
	serverName    string
	serverVersion string
	handler       modelcontextprotocol.McpServerEventHandler
	// used by protocol
	clientName          string
	clientVersion       string
	isClientInitialized bool
	lastRequestId       int
	jsonRpcTransport    *transport.JsonRpcTransport
}

// constructor for the MCP server
func NewMcpSdkServer(serverDefinition types.McpSdkServerDefinition, debug bool) (types.ModelContextProtocolServer, error) {
	// We get the concrete type of the server definition
	sdkServerDefinition, ok := serverDefinition.(*sdk.SdkServerDefinition)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type: expected *sdk.SdkServerDefinition, got %T", serverDefinition)
	}

	// we build the configuration data
	loggingInfo := &logger.LoggingInfo{
		Level:      sdkServerDefinition.DebugLevel(),
		File:       sdkServerDefinition.DebugFile(),
		WithStderr: false,
	}

	// we initialize the logger
	logger, err := logger.NewLogger(loggingInfo, debug)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	// we create the MCP server handler
	mcpServerNotifications, err := providers.NewProviderMcpServerHandler(sdkServerDefinition, logger)
	if err != nil {
		return nil, err
	}

	return &McpServer{
		logger:        logger,
		serverName:    sdkServerDefinition.ServerName(),
		serverVersion: sdkServerDefinition.ServerVersion(),
		handler:       mcpServerNotifications,
		lastRequestId: 0,
	}, nil

}

func (mcp *McpServer) StdioTransport() types.Transport {
	// we create the transport
	transport := transport.NewStdioTransport(
		mcp.logger)

	// we return the transport
	return transport
}
