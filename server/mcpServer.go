package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/jsonrpc"
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
	logger              types.Logger
}

func NewMCPServer(
	transport types.Transport,
	toolsRegistry *tools.ToolsRegistry,
	promptsRegistry *prompts.PromptsRegistry,
	serverName string,
	serverVersion string,
	logger types.Logger,
) *MCPServer {
	return &MCPServer{
		transport:       transport,
		toolsRegistry:   toolsRegistry,
		promptsRegistry: promptsRegistry,
		serverName:      serverName,
		serverVersion:   serverVersion,
		logger:          logger,
	}
}

func (s *MCPServer) Start(ctx context.Context) error {
	transport := s.transport

	// Set up message handler
	transport.OnMessage(func(msg json.RawMessage) {
		nature, jsonRpcRawMessage, err := jsonrpc.CheckJsonMessage(msg)

		if err != nil || nature != jsonrpc.MessageNatureRequest {
			s.logger.Debug("invalid message received message", types.LogArg{
				"message": string(msg),
			})
			s.sendError(&jsonrpc.JsonRpcError{
				Code:    jsonrpc.RpcParseError,
				Message: "invalid message",
			}, nil)
		}

		request, requestId, rpcErr := jsonrpc.ParseJsonRpcRequest(jsonRpcRawMessage)
		if rpcErr != nil {
			s.sendError(rpcErr, requestId)
			return
		}

		err = s.processRequest(ctx, request)
		if err != nil {
			s.logError("failed to process request", err)
		}

	})

	// Set up error handler
	transport.OnError(func(err error) {
		s.logError("transport error", err)
	})

	// Start the transport
	if err := transport.Start(ctx); err != nil {
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
	s.logger.Error(message, types.LogArg{
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
	s.logger.Debug("JsonRpcError", types.LogArg{
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
