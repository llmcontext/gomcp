package mcp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/inspector"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/mux"
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
	multiplexer     *mux.Multiplexer
	logger          types.Logger
}

func NewModelContextProtocolServer(configFilePath string) (*ModelContextProtocolImpl, error) {
	// we load the config file
	config, err := config.LoadConfig(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %v", configFilePath, err)
	}

	// we initialize the logger
	logger, err := logger.NewLogger(config.Logging, false)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	// Initialize tools registry
	toolsRegistry := tools.NewToolsRegistry(logger)

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
		inspectorInstance = inspector.NewInspector(config.Inspector, logger)
	}

	// Start multiplexer if enabled
	var multiplexerInstance *mux.Multiplexer = nil
	if config.Proxy != nil && config.Proxy.Enabled {
		multiplexerInstance = mux.NewMultiplexer(config.Proxy, logger)
	}

	return &ModelContextProtocolImpl{
		config:          config,
		toolsRegistry:   toolsRegistry,
		promptsRegistry: promptsRegistry,
		inspector:       inspectorInstance,
		multiplexer:     multiplexerInstance,
		logger:          logger,
	}, nil
}

func (mcp *ModelContextProtocolImpl) StdioTransport() types.Transport {
	// delete the protocol debug file if it exists
	if mcp.config.Logging.ProtocolDebugFile != "" {
		if _, err := os.Stat(mcp.config.Logging.ProtocolDebugFile); err == nil {
			os.Remove(mcp.config.Logging.ProtocolDebugFile)
		}
	}

	// we create the transport
	transport := transport.NewStdioTransport(
		mcp.config.Logging.ProtocolDebugFile,
		mcp.inspector,
		mcp.logger)

	// we return the transport
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

	// Start inspector if it was enabled
	if mcp.inspector != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mcp.inspector.Start(ctx)
		}()
	}

	// Start MCP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Initialize server
		server := server.NewMCPServer(transport, mcp.toolsRegistry, mcp.promptsRegistry,
			mcp.config.ServerInfo.Name,
			mcp.config.ServerInfo.Version,
			mcp.logger)

		// Start server
		err = server.Start(ctx)
		if err != nil {
			mcp.logger.Error("error starting server", types.LogArg{
				"error": err,
			})
			cancel()
		}
	}()

	// Start multiplexer if it was enabled
	if mcp.multiplexer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := mcp.multiplexer.Start(ctx)
			if err != nil {
				mcp.logger.Error("error starting multiplexer", types.LogArg{
					"error": err,
				})
				cancel()
			}
		}()
	}

	// Listen for OS signals (e.g., Ctrl+C)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGCHLD)

	go func() {
		for {
			parentPID := syscall.Getppid()
			mcp.logger.Info("Monitoring parent process", types.LogArg{
				"pid": parentPID,
			})
			if parentPID == 1 {
				mcp.logger.Info("Parent process is init. Shutting down...", types.LogArg{
					"pid": parentPID,
				})
				signalChan <- os.Interrupt
				return
			}
			// we wait 10 seconds before checking again
			time.Sleep(10 * time.Second)
		}
	}()

	// Wait for a signal to stop the server
	sig := <-signalChan
	fmt.Fprintf(os.Stderr, "[mcp] Received an interrupt, shutting down... %s\n", sig)

	// Cancel the context to signal the goroutines to stop
	cancel()

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Fprintf(os.Stderr, "[mcp] All goroutines have stopped. Exiting.\n")

	return nil
}

func (mcp *ModelContextProtocolImpl) GetToolRegistry() types.ToolRegistry {
	return mcp
}
