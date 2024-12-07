package gomcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/server"
	"github.com/llmcontext/gomcp/tools"
	"github.com/llmcontext/gomcp/transport"
)

type ToolProvider interface {
	AddTool(toolName string, description string, toolHandler interface{}) error
}

type ModelContextProtocol struct {
	toolsRegistry *tools.ToolsRegistry
	config        *config.Config
}

func NewModelContextProtocolServer(configFilePath string) (*ModelContextProtocol, error) {
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

	return &ModelContextProtocol{
		toolsRegistry: toolsRegistry,
		config:        config,
	}, nil
}

func (mcp *ModelContextProtocol) CreateStdioTransport() transport.Transport {
	transport := transport.NewStdioTransport()
	return transport
}

func (mcp *ModelContextProtocol) DeclareToolProvider(toolName string, toolInitFunction interface{}) (ToolProvider, error) {
	toolProvider, err := tools.DeclareToolProvider(toolName, toolInitFunction)
	if err != nil {
		return nil, fmt.Errorf("failed to declare tool provider %s: %v", toolName, err)
	}
	// we keep track of the tool providers added
	mcp.toolsRegistry.RegisterToolProvider(toolProvider)
	return toolProvider, nil
}

func (mcp *ModelContextProtocol) Start(serverName string, serverVersion string, transport transport.Transport) error {
	// All the tools are initialized, we can prepare the tools registry
	// so that it can be used by the server
	err := mcp.toolsRegistry.Prepare(mcp.config.Tools)
	if err != nil {
		return fmt.Errorf("error preparing tools registry: %s", err)
	}

	// Initialize server
	server := server.NewMCPServer(transport, mcp.toolsRegistry, serverName, serverVersion)

	// Start server
	err = server.Start()
	if err != nil {
		return fmt.Errorf("error starting server: %s", err)
	}
	return nil
}
