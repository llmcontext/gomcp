package mcpserver

import (
	"fmt"

	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/modelcontextprotocol"
	"github.com/llmcontext/gomcp/providers"
	"github.com/llmcontext/gomcp/providers/presets"
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

// constructor for the MCP server developed with the SDK
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

	return newMcpServer(
		logger,
		sdkServerDefinition.ServerName(),
		sdkServerDefinition.ServerVersion(),
		mcpServerNotifications,
	), nil
}

// constructor for the multiplexer MCP server
func NewMcpServer(serverName string, serverVersion string, loggingInfo *logger.LoggingInfo, debug bool) (*McpServer, error) {
	logger, err := logger.NewLogger(loggingInfo, debug)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	logger.Debug("registry>server>NewMcpServer", types.LogArg{
		"serverName":    serverName,
		"serverVersion": serverVersion,
	})

	// Register preset servers
	// we use the same registration mechanism as for the SDK servers
	sdkServerDefinition := sdk.NewMcpSdkServerDefinition(serverName, serverVersion)
	presets.RegisterPresetServers(sdkServerDefinition, logger)

	// we create the MCP server handler
	mcpServerHandler, err := providers.NewProviderMcpServerHandler(sdkServerDefinition, logger)
	if err != nil {
		return nil, err
	}

	return newMcpServer(
		logger,
		serverName,
		serverVersion,
		mcpServerHandler,
	), nil

}

// common constructor for the MCP server
func newMcpServer(
	logger types.Logger,
	serverName string,
	serverVersion string,
	handler modelcontextprotocol.McpServerEventHandler,
) *McpServer {

	return &McpServer{
		logger:        logger,
		serverName:    serverName,
		serverVersion: serverVersion,
		handler:       handler,
		lastRequestId: 0,
	}
}

func (mcp *McpServer) StdioTransport() types.Transport {
	// we create the transport
	transport := transport.NewStdioTransport(
		mcp.logger)

	// we return the transport
	return transport
}
