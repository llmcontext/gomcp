package mcp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/inspector"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/prompts"
	"github.com/llmcontext/gomcp/server"
	"github.com/llmcontext/gomcp/tools"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type ModelContextProtocolImpl struct {
	config          *config.Config
	toolsRegistry   *tools.ToolsRegistry
	promptsRegistry *prompts.PromptsRegistry
	inspector       *inspector.Inspector
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

	// Initialize prompts registry
	promptsRegistry := prompts.NewEmptyPromptsRegistry()
	if config.Prompts != nil {
		promptsRegistry, err = prompts.NewPromptsRegistry(config.Prompts.File)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize prompts registry: %v", err)
		}
	}

	// Start inspector if enabled
	var inspectorInstance *inspector.Inspector = nil
	if config.Inspector != nil && config.Inspector.Enabled {
		inspectorInstance = inspector.NewInspector(config.Inspector)
	}

	return &ModelContextProtocolImpl{
		config:          config,
		toolsRegistry:   toolsRegistry,
		promptsRegistry: promptsRegistry,
		inspector:       inspectorInstance,
	}, nil
}

func (mcp *ModelContextProtocolImpl) StdioTransport() types.Transport {
	transport := transport.NewStdioTransport(
		mcp.config.Logging.ProtocolDebugFile,
		mcp.inspector)
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

// Start starts the server and the inspector
func (mcp *ModelContextProtocolImpl) Start(transport types.Transport) error {
	// we create a context that will be used to cancel the server and the inspector
	ctx, cancel := context.WithCancel(context.Background())

	// Use a wait group to wait for goroutines to complete
	var wg sync.WaitGroup

	// All the tools are initialized, we can prepare the tools registry
	// so that it can be used by the server
	err := mcp.toolsRegistry.Prepare(ctx, mcp.config.Tools)
	if err != nil {
		cancel()
		return fmt.Errorf("error preparing tools registry: %s", err)
	}

	// Start inspector if enabled
	if mcp.config.Inspector != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mcp.inspector.StartInspector(ctx)
		}()
	}

	// Start MCP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Initialize server
		server := server.NewMCPServer(transport, mcp.toolsRegistry, mcp.promptsRegistry,
			mcp.config.ServerInfo.Name,
			mcp.config.ServerInfo.Version)

		// Start server
		err = server.Start(ctx)
		if err != nil {
			logger.Error("error starting server", logger.Arg{
				"error": err,
			})
			cancel()
		}
	}()

	// Listen for OS signals (e.g., Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Wait for a signal to stop the server
	<-signalChan
	fmt.Println("\nReceived an interrupt, shutting down...")

	// Cancel the context to signal the goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All goroutines have stopped. Exiting.")

	return nil
}

func (mcp *ModelContextProtocolImpl) GetToolRegistry() types.ToolRegistry {
	return mcp
}
