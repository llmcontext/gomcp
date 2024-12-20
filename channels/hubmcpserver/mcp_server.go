package hubmcpserver

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/llmcontext/gomcp/channels/hub/events"
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
	events          events.Events
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
	events events.Events,
	toolsRegistry *tools.ToolsRegistry,
	promptsRegistry *prompts.PromptsRegistry,
	serverName string,
	serverVersion string,
	logger types.Logger,
) *MCPServer {
	return &MCPServer{
		transport:       transport,
		events:          events,
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

	errChan := make(chan error, 1)

	go func() {
		// Start the transport
		err := transport.Start(ctx)
		if err != nil {
			s.logError("failed to start transport", err)
		}
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		transport.Close()
		return ctx.Err()
	}
}

func (s *MCPServer) Close() {
	s.transport.Close()
}

func (s *MCPServer) logError(message string, err error) {
	// check if the error is because the context was cancelled
	if errors.Is(err, context.Canceled) {
		s.logger.Info("transport - context cancelled", types.LogArg{})
	} else {
		s.logger.Error(message, types.LogArg{
			"message": message,
			"error":   err,
		})
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

func (s *MCPServer) OnNewProxyTools() {
	// TODO: implement
	tools := s.toolsRegistry.GetListOfTools()
	s.logger.Info("OnNewProxyTools", types.LogArg{
		"tools": tools,
	})

}
