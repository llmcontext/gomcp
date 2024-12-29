package hub

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/llmcontext/gomcp/channels/hub/events"
	"github.com/llmcontext/gomcp/channels/hubinspector"
	"github.com/llmcontext/gomcp/channels/hubmcpserver"
	"github.com/llmcontext/gomcp/channels/hubmuxserver"
	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/prompts"
	"github.com/llmcontext/gomcp/sdk"
	"github.com/llmcontext/gomcp/servers/mcp_server_time"
	"github.com/llmcontext/gomcp/tools"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
	"golang.org/x/sync/errgroup"
)

type ModelContextProtocolImpl struct {
	logger          types.Logger
	toolsRegistry   *tools.ToolsRegistry
	promptsRegistry *prompts.PromptsRegistry
	inspector       *hubinspector.Inspector
	muxServer       *hubmuxserver.MuxServer
	stateManager    *StateManager
	events          events.Events
}

func newModelContextProtocolServer(
	serverInfo *config.ServerInfo,
	logger types.Logger,
	promptsConfig *config.PromptConfig,
	inspectorConfig *config.InspectorInfo,
	toolsRegistry *tools.ToolsRegistry,
	proxyConfig *config.ServerProxyConfig) (*ModelContextProtocolImpl, error) {
	// we initialize the logger
	// logger, err := logger.NewLogger(logging, false)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to initialize logger: %v", err)
	// }
	var err error

	// Initialize prompts registry
	promptsRegistry := prompts.NewEmptyPromptsRegistry()
	if promptsConfig != nil {
		promptsRegistry, err = prompts.NewPromptsRegistry(promptsConfig.File)
		if err != nil {
			return nil, fmt.Errorf("failed to initialize prompts registry: %v", err)
		}
	}

	// initialize the state manager
	stateManager := NewStateManager(
		serverInfo.Name,
		serverInfo.Version,
		toolsRegistry,
		promptsRegistry,
		logger,
	)
	events := stateManager.AsEvents()

	// Start inspector if enabled
	var inspectorInstance *hubinspector.Inspector = nil
	if inspectorConfig != nil && inspectorConfig.Enabled {
		inspectorInstance = hubinspector.NewInspector(inspectorConfig, logger)
	}

	// Start multiplexer if enabled
	var muxServerInstance *hubmuxserver.MuxServer = nil
	if proxyConfig != nil && proxyConfig.Enabled {
		muxServerInstance = hubmuxserver.NewMuxServer(proxyConfig.ListenAddress, events, logger)
	}

	return &ModelContextProtocolImpl{
		logger:          logger,
		toolsRegistry:   toolsRegistry,
		promptsRegistry: promptsRegistry,
		inspector:       inspectorInstance,
		muxServer:       muxServerInstance,
		stateManager:    stateManager,
		events:          events,
	}, nil

}

func NewHubModelContextProtocolServer(debug bool) (*ModelContextProtocolImpl, error) {
	conf, err := config.LoadHubConfiguration()
	if err != nil {
		return nil, fmt.Errorf("failed to load hub configuration: %v", err)
	}

	if debug {
		conf.Logging.WithStderr = true
	}

	// we initialize the logger
	logger, err := logger.NewLogger(conf.Logging, debug)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	toolsRegistry := tools.NewToolsRegistry(true, logger)

	return newModelContextProtocolServer(
		&conf.ServerInfo,
		logger,
		conf.Prompts,
		conf.Inspector,
		toolsRegistry,
		conf.Proxy,
	)
}

func (mcp *ModelContextProtocolImpl) LoadPresetTools() error {
	presetToolsNames := []string{
		"gomcp_server_time",
	}

	// TODO: add mechanism to disable some preset tools
	for _, toolName := range presetToolsNames {
		switch toolName {
		case "gomcp_server_time":
			err := mcp_server_time.RegisterTools(mcp.toolsRegistry)
			if err != nil {
				mcp.logger.Error("failed to register tools: %v", types.LogArg{
					"error": err,
				})
			}
		}
	}
	return nil
}

func NewModelContextProtocolServer(configuration types.McpServerDefinition) (*ModelContextProtocolImpl, error) {
	// Add type assertion check
	sdkConfiguration, ok := configuration.(*sdk.SdkServerDefinition)
	if !ok {
		return nil, fmt.Errorf("invalid configuration type: expected *sdk.SdkServerDefinition, got %T", configuration)
	}

	// create the McpServerRegistry
	mcpServerRegistry := NewMcpServerRegistry()

	// we build the configuration data
	conf := config.ServerConfiguration{
		ServerInfo: config.ServerInfo{
			Name:    sdkConfiguration.ServerName(),
			Version: sdkConfiguration.ServerVersion(),
		},
		Logging: &config.LoggingInfo{
			Level:      sdkConfiguration.DebugLevel(),
			File:       sdkConfiguration.DebugFile(),
			WithStderr: false,
		},
	}

	// we initialize the logger
	logger, err := logger.NewLogger(conf.Logging, false)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	toolsRegistry := tools.NewToolsRegistry(false, logger)

	// Setup the SDK based MCP servers
	err = sdkConfiguration.SetupMcpServer(toolsRegistry, mcpServerRegistry)
	if err != nil {
		return nil, err
	}

	return newModelContextProtocolServer(
		&conf.ServerInfo,
		logger,
		conf.Prompts,
		conf.Inspector,
		toolsRegistry,
		nil,
	)

}

func (mcp *ModelContextProtocolImpl) StdioTransport() types.Transport {
	// we create the transport
	transport := transport.NewStdioTransport(
		mcp.inspector,
		mcp.logger)

	// we return the transport
	return transport
}

// Start starts the server and the inspector
func (mcp *ModelContextProtocolImpl) Start(transport types.Transport) error {
	mcp.logger.Info("Starting MCP server", types.LogArg{})

	// create a context that will be used to cancel the server and the inspector
	ctx := context.Background()

	// All the tools are initialized, we can prepare the tools registry
	// so that it can be used by the server
	err := mcp.toolsRegistry.Prepare(ctx)
	if err != nil {
		return fmt.Errorf("error preparing tools registry: %s", err)
	}

	mcp.logger.Info("Starting inspector", types.LogArg{})

	// we create an errgroup that will be used to cancel
	// all the components of the server
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		mcp.logger.Info("[A] for MCP server to stop", types.LogArg{})

		// Listen for OS signals (e.g., Ctrl+C)
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGCHLD)

		select {
		case <-egCtx.Done():
			mcp.logger.Info("[A.1] Context cancelled, shutting down", types.LogArg{})
			return egCtx.Err()
		case signal := <-signalChan:
			mcp.logger.Info("[A.2] Received an interrupt, shutting down", types.LogArg{
				"signal": signal.String(),
			})
			return fmt.Errorf("received an interrupt (%s)", signal.String())
		}
	})

	// Start inspector if it was enabled
	if mcp.inspector != nil {
		eg.Go(func() error {
			mcp.logger.Info("[B] Starting inspector", types.LogArg{})
			err := mcp.inspector.Start(egCtx)
			if err != nil {
				// check if the error is because the context was cancelled
				if errors.Is(err, context.Canceled) {
					mcp.logger.Info("[B.1] context cancelled, stopping inspector", types.LogArg{})
				} else {
					mcp.logger.Error("[B.2] error starting inspector", types.LogArg{
						"error": err,
					})
				}
			}
			mcp.logger.Info("[B.3] inspector stopped", types.LogArg{})
			return err
		})
	}

	eg.Go(func() error {
		mcp.logger.Info("[C] Starting MCP server", types.LogArg{})

		// Initialize server
		server := hubmcpserver.NewMCPServer(transport,
			mcp.events,
			mcp.logger)

		// set the server in the state manager
		mcp.stateManager.SetMcpServer(server)

		// Start server
		err := server.Start(egCtx)
		if err != nil {
			// check if the error is because the context was cancelled
			if errors.Is(err, context.Canceled) {
				mcp.logger.Info("[C.1] context cancelled, stopping MCP server", types.LogArg{})
			} else {
				mcp.logger.Error("[C.2] error starting MCP server", types.LogArg{
					"error": err,
				})
			}
		}
		mcp.logger.Info("[C.3] MCP server stopped", types.LogArg{})
		return err
	})

	// Start multiplexer if it was enabled
	if mcp.muxServer != nil {
		eg.Go(func() error {
			mcp.logger.Info("[D] Starting mux server", types.LogArg{})
			mcp.stateManager.SetMuxServer(mcp.muxServer)

			err := mcp.muxServer.Start(egCtx)
			if err != nil {
				// check if the error is because the context was cancelled
				if errors.Is(err, context.Canceled) {
					mcp.logger.Info("[D.1] context cancelled, stopping multiplexer", types.LogArg{})
				} else {
					mcp.logger.Error("[D.2] error starting multiplexer", types.LogArg{
						"error": err,
					})
				}
			}
			mcp.logger.Info("[D.3] mux server stopped", types.LogArg{})
			return err
		})
	}

	if false {
		eg.Go(func() error {
			count := 0
			timer := time.NewTimer(10 * time.Second)
			defer timer.Stop()

			for {
				count++
				parentPID := syscall.Getppid()
				mcp.logger.Info("Monitoring parent process", types.LogArg{
					"pid":   parentPID,
					"count": count,
				})
				if parentPID == 1 {
					mcp.logger.Info("Parent process is init. Shutting down...", types.LogArg{
						"pid": parentPID,
					})
					return fmt.Errorf("parent process is init")
				}
				logGoroutineStacks(mcp.logger)

				if count > 10 {
					mcp.logger.Info("Stopping parent process monitor", types.LogArg{})
					return nil
				}

				timer.Reset(10 * time.Second)
				select {
				case <-egCtx.Done():
					mcp.logger.Info("Context cancelled, stopping parent process monitor", types.LogArg{})
					return egCtx.Err()
				case <-timer.C:
					continue
				}
			}
		})
	}

	err = eg.Wait()
	if err != nil {
		mcp.logger.Info("Stopping hub server", types.LogArg{
			"reason": err.Error(),
		})
	}
	mcp.logger.Info("Hub server stopped", types.LogArg{})
	return nil
}

func logGoroutineStacks(logger types.Logger) {
	// Get number of goroutines
	numGoroutines := runtime.NumGoroutine()

	// Get stack traces
	buf := make([]byte, 1024*10)
	n := runtime.Stack(buf, true)
	stacks := string(buf[:n])

	logger.Info("Goroutine dump", types.LogArg{
		"num_goroutines": numGoroutines,
	})

	// print the goroutine stacks formatted
	fmt.Println("Goroutine stacks:")
	fmt.Println(stacks)
}
