package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/logger"
	"github.com/llmcontext/gomcp/prompts"
	"github.com/llmcontext/gomcp/tools"
	"github.com/llmcontext/gomcp/types"
)

type ClientInfo struct {
	name    string
	version string
}

type MCPServer struct {
	transport       types.Transport
	toolsRegistry   *tools.ToolsRegistry
	promptsRegistry *prompts.PromptsRegistry
	// server information
	serverName    string
	serverVersion string
	// client information
	isClientInitialized bool
	protocolVersion     string
	clientInfo          *ClientInfo
}

func NewMCPServer(transport types.Transport, toolsRegistry *tools.ToolsRegistry, promptsRegistry *prompts.PromptsRegistry, serverName string, serverVersion string) *MCPServer {
	return &MCPServer{
		transport:       transport,
		toolsRegistry:   toolsRegistry,
		promptsRegistry: promptsRegistry,
		serverName:      serverName,
		serverVersion:   serverVersion,
	}
}

func (s *MCPServer) Start(ctx context.Context) error {
	transport := s.transport

	// Set up message handler
	transport.OnMessage(func(msg json.RawMessage) {
		requests, isBatch, error := jsonrpc.ParseRequest(msg)
		if error != nil {
			logger.Debug("invalid transport received message", logger.Arg{
				"message": string(msg),
			})
			s.sendError(error, nil)
			return
		}
		if isBatch {
			s.logError("batched requests not supported yet", nil)
		} else {
			request := requests[0]
			if request.Error != nil {
				s.sendError(request.Error, request.RequestId)
				return
			}
			err := s.processRequest(ctx, request.Request)
			if err != nil {
				s.logError("failed to process request", err)
			}
		}
	})

	// Set up error handler
	transport.OnError(func(err error) {
		s.logError("transport error", err)
	})

	// Start the transport
	if err := transport.Start(); err != nil {
		s.logError("failed to start transport", err)
		return err
	}

	// Keep the main thread alive
	// will be interrupted by the context
	<-ctx.Done()

	fmt.Printf("# [mcpServer] shutdown\n")

	transport.Close()
	return nil
}

func (s *MCPServer) logError(message string, err error) {
	logger.Error(message, logger.Arg{
		"message": message,
		"error":   err,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", message)
	}
}

func (s *MCPServer) sendError(error *jsonrpc.JsonRpcError, id *jsonrpc.JsonRpcRequestId) {
	logger.Debug("JsonRpcError", logger.Arg{
		"error": error,
		"id":    id,
	})
	response := &jsonrpc.JsonRpcResponse{
		Error: error,
		Id:    id,
	}
	jsonError, err := jsonrpc.MarshalJsonRpcResponse(response)
	if err != nil {
		s.logError("failed to marshal error", err)
		return
	}
	s.transport.Send(jsonError)
}
