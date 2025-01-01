package mcpserver

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/llmcontext/gomcp/types"
	"golang.org/x/sync/errgroup"
)

func (m *McpServer) Start(transport types.Transport) error {
	var err error
	m.logger.Info("Starting MCP server", types.LogArg{})

	// create a context that will be used to cancel the server and the inspector
	ctx := context.Background()

	// we create an errgroup that will be used to cancel
	// all the components of the server
	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		// Listen for OS signals (e.g., Ctrl+C)
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGCHLD)

		select {
		case <-egCtx.Done():
			m.logger.Info("Context cancelled, shutting down", types.LogArg{})
			return egCtx.Err()
		case signal := <-signalChan:
			m.logger.Info("Received an interrupt, shutting down", types.LogArg{
				"signal": signal.String(),
			})
			return fmt.Errorf("received an interrupt (%s)", signal.String())
		}
	})

	// Start inspector if it was enabled
	// if mcp.inspector != nil {
	// 	eg.Go(func() error {
	// 		mcp.logger.Info("[B] Starting inspector", types.LogArg{})
	// 		err := mcp.inspector.Start(egCtx)
	// 		if err != nil {
	// 			// check if the error is because the context was cancelled
	// 			if errors.Is(err, context.Canceled) {
	// 				mcp.logger.Info("[B.1] context cancelled, stopping inspector", types.LogArg{})
	// 			} else {
	// 				mcp.logger.Error("[B.2] error starting inspector", types.LogArg{
	// 					"error": err,
	// 				})
	// 			}
	// 		}
	// 		mcp.logger.Info("[B.3] inspector stopped", types.LogArg{})
	// 		return err
	// 	})
	// }

	eg.Go(func() error {
		m.logger.Info("Starting MCP protocol", types.LogArg{})

		// Start protocol management
		err := m.startProtocol(egCtx, transport)
		if err != nil {
			// check if the error is because the context was cancelled
			if errors.Is(err, context.Canceled) {
				m.logger.Info("context cancelled, stopping MCP server", types.LogArg{})
			} else {
				m.logger.Error("error starting MCP server", types.LogArg{
					"error": err,
				})
			}
		}
		m.logger.Info("MCP transport stopped", types.LogArg{})
		return err
	})

	err = eg.Wait()
	if err != nil {
		m.logger.Info("Stopping MCP server", types.LogArg{
			"reason": err.Error(),
		})
	}
	m.logger.Info("MCP server stopped", types.LogArg{})
	return nil
}
