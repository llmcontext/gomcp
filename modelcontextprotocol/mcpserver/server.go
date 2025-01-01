package mcpserver

import (
	"fmt"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/modelcontextprotocol"
	"github.com/llmcontext/gomcp/providers"
	"github.com/llmcontext/gomcp/providers/presets"
	"github.com/llmcontext/gomcp/providers/sdk"
	"github.com/llmcontext/gomcp/registry"
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
	loggingInfo := &config.LoggingInfo{
		Level:      sdkServerDefinition.DebugLevel(),
		File:       sdkServerDefinition.DebugFile(),
		WithStderr: false,
	}

	// we initialize the logger
	logger, err := logger.NewLogger(loggingInfo, debug)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	// create the McpServerRegistry
	mcpServerRegistry := registry.NewMcpServerRegistry(logger)

	// Setup the SDK based MCP servers
	err = sdkServerDefinition.RegisterSdkMcpServer(mcpServerRegistry)
	if err != nil {
		return nil, err
	}

	mcpServerNotifications, err := providers.NewProviderMcpServerHandler(sdkServerDefinition, false, logger)
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
func NewMcpServer(serverInfo *config.ServerInfo, loggingInfo *config.LoggingInfo, debug bool) (*McpServer, error) {
	logger, err := logger.NewLogger(loggingInfo, debug)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	// create the McpServerRegistry
	mcpServerRegistry := registry.NewMcpServerRegistry(logger)

	logger.Debug("registry>server>NewMcpServer", types.LogArg{
		"serverName":    serverInfo.Name,
		"serverVersion": serverInfo.Version,
	})

	// Register preset servers
	// we use the same registration mechanism as for the SDK servers
	serverDefinition := sdk.NewMcpServerDefinition(serverInfo.Name, serverInfo.Version)
	presets.RegisterPresetServers(serverDefinition, logger)
	sdkServerDefinition, ok := serverDefinition.(*sdk.SdkServerDefinition)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type: expected *sdk.SdkServerDefinition, got %T", serverDefinition)
	}

	// Setup the SDK based MCP servers
	err = sdkServerDefinition.RegisterSdkMcpServer(mcpServerRegistry)
	if err != nil {
		return nil, err
	}

	mcpServerHandler, err := providers.NewProviderMcpServerHandler(sdkServerDefinition, true, logger)
	if err != nil {
		return nil, err
	}

	return newMcpServer(
		logger,
		serverInfo.Name,
		serverInfo.Version,
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
