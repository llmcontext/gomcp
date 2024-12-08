package mcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/server"
	"github.com/llmcontext/gomcp/tools"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type ModelContextProtocolImpl struct {
	toolsRegistry *tools.ToolsRegistry
	config        *config.Config
}

func NewModelContextProtocolServer(configFilePath string) (*ModelContextProtocolImpl, error) {
	// we load the config file
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %v", configFilePath, err)
	}

	// we initialize the logger
	err = logger.InitLogger(config.Logging, false)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	// Initialize tools registry
	toolsRegistry := tools.NewToolsRegistry()

	return &ModelContextProtocolImpl{
		toolsRegistry: toolsRegistry,
		config:        config,
	}, nil
}

func (mcp *ModelContextProtocolImpl) StdioTransport() types.Transport {
	transport := transport.NewStdioTransport()
	return transport
}

func (mcp *ModelContextProtocolImpl) DeclareToolProvider(toolName string, toolInitFunction interface{}) (types.ToolProvider, error) {
	toolProvider, err := tools.DeclareToolProvider(toolName, toolInitFunction)
	if err != nil {
		return nil, fmt.Errorf("failed to declare tool provider %s: %v", toolName, err)
	}
	// we keep track of the tool providers added
	mcp.toolsRegistry.RegisterToolProvider(toolProvider)
	return toolProvider, nil
}

func (mcp *ModelContextProtocolImpl) Start(transport types.Transport) error {
	// All the tools are initialized, we can prepare the tools registry
	// so that it can be used by the server
	err := mcp.toolsRegistry.Prepare(mcp.config.Tools)
	if err != nil {
		return fmt.Errorf("error preparing tools registry: %s", err)
	}

	// Initialize server
	server := server.NewMCPServer(transport, mcp.toolsRegistry,
		mcp.config.ServerInfo.Name,
		mcp.config.ServerInfo.Version)

	// Start server
	err = server.Start()
	if err != nil {
		return fmt.Errorf("error starting server: %s", err)
	}
	return nil
}

func (mcp *ModelContextProtocolImpl) GetToolRegistry() types.ToolRegistry {
	return mcp
}
