package mcp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/endpoints/mcpserver"
	"github.com/llmcontext/gomcp/endpoints/muxserver"
	"github.com/llmcontext/gomcp/eventbus"
	"github.com/llmcontext/gomcp/inspector"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/prompts"
	"github.com/llmcontext/gomcp/tools"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
	"golang.org/x/sync/errgroup"
)

type ModelContextProtocolImpl struct {
	config          *config.Config
	eventBus        *eventbus.EventBus
	toolsRegistry   *tools.ToolsRegistry
	promptsRegistry *prompts.PromptsRegistry
	inspector       *inspector.Inspector
	multiplexer     *muxserver.Multiplexer
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

	// Initialize the event bus used to communicate between the
	// different components of the server
	eventBus := eventbus.NewEventBus()

	// Initialize tools registry
	toolsRegistry := tools.NewToolsRegistry(eventBus, logger)

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
	var multiplexerInstance *muxserver.Multiplexer = nil
	if config.Proxy != nil && config.Proxy.Enabled {
		multiplexerInstance = muxserver.NewMultiplexer(config.Proxy, eventBus, logger)
	}

	return &ModelContextProtocolImpl{
		config:          config,
		eventBus:        eventBus,
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
	mcp.logger.Info("Starting MCP server", types.LogArg{})

	// create a context that will be used to cancel the server and the inspector
	ctx := context.Background()

	// All the tools are initialized, we can prepare the tools registry
	// so that it can be used by the server
	err := mcp.toolsRegistry.Prepare(ctx, mcp.config.Tools)
	if err != nil {
		return fmt.Errorf("error preparing tools registry: %s", err)
	}

	mcp.logger.Info("Starting inspector", types.LogArg{})

	// we create an errgroup that will be used to cancel the server and the inspector
	eg, egCtx := errgroup.WithContext(ctx)

	// Start inspector if it was enabled
	if mcp.inspector != nil {
		eg.Go(func() error {
			err := mcp.inspector.Start(egCtx)
			if err != nil {
				mcp.logger.Error("error starting inspector", types.LogArg{
					"error": err,
				})
			}
			return nil
		})
	}

	eg.Go(func() error {
		mcp.logger.Info("Starting MCP server", types.LogArg{})

		// Initialize server
		server := mcpserver.NewMCPServer(transport, mcp.toolsRegistry, mcp.promptsRegistry,
			mcp.config.ServerInfo.Name,
			mcp.config.ServerInfo.Version,
			mcp.logger)

		// Start server
		err := server.Start(egCtx)
		if err != nil {
			mcp.logger.Error("error starting server", types.LogArg{
				"error": err,
			})
			return err
		}
		return nil
	})

	// Start multiplexer if it was enabled
	if mcp.multiplexer != nil {
		eg.Go(func() error {
			mcp.logger.Info("Starting multiplexer", types.LogArg{})
			err := mcp.multiplexer.Start(egCtx)
			if err != nil {
				mcp.logger.Error("error starting multiplexer", types.LogArg{
					"error": err,
				})
				return err
			}
			return nil
		})
	}

	eg.Go(func() error {
		mcp.logger.Info("Waiting for MCP server to stop", types.LogArg{})

		// Listen for OS signals (e.g., Ctrl+C)
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGCHLD)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case signal := <-signalChan:
			mcp.logger.Info("Received an interrupt, shutting down...", types.LogArg{
				"signal": signal.String(),
			})
			return fmt.Errorf("received an interrupt (%s)", signal.String())
		}
	})

	eg.Go(func() error {
		timer := time.NewTimer(10 * time.Second)
		defer timer.Stop()

		for {
			parentPID := syscall.Getppid()
			// mcp.logger.Info("Monitoring parent process", types.LogArg{
			// 	"pid": parentPID,
			// })
			if parentPID == 1 {
				mcp.logger.Info("Parent process is init. Shutting down...", types.LogArg{
					"pid": parentPID,
				})
				return fmt.Errorf("parent process is init")
			}

			timer.Reset(10 * time.Second)
			select {
			case <-ctx.Done():
				mcp.logger.Info("Context cancelled, stopping parent process monitor", types.LogArg{})
				return ctx.Err()
			case <-timer.C:
				continue
			}
		}
	})

	err = eg.Wait()
	if err != nil {
		mcp.logger.Error("stopping server", types.LogArg{
			"error": err,
		})
	}
	return nil
}

func (mcp *ModelContextProtocolImpl) GetToolRegistry() types.ToolRegistry {
	return mcp
}
