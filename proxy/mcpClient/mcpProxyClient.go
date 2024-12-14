package mcpClient

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/proxy/muxClient"
	"github.com/llmcontext/gomcp/types"
	"github.com/llmcontext/gomcp/version"
)

type MCPProxyClientOptions struct {
	ProxyName               string
	CurrentWorkingDirectory string
	ProgramName             string
	ProgramArgs             []string
}

type PendingRequest struct {
	method    string
	messageId *jsonrpc.JsonRpcRequestId
}

type MCPProxyClient struct {
	options MCPProxyClientOptions
	logger  types.TermLogger
	// transport is the transport layer to communicate with the actual MCP server
	transport types.Transport
	// muxClient is the client to communicate with the mux server
	muxClient *muxClient.MuxClient
	// pendingRequests is a map of message id to pending request
	pendingRequests map[string]*PendingRequest
	// clientId is the id of the next message to send
	clientId int

	// serverInfo is the info about the server we are connected to
	serverInfo mcp.ServerInfo
	// tools is the list of tools available for the proxy
	tools []mcp.ToolDescription
}

func NewMCPProxyClient(transport types.Transport, muxClient *muxClient.MuxClient, options MCPProxyClientOptions, logger types.TermLogger) *MCPProxyClient {
	return &MCPProxyClient{
		transport:       transport,
		muxClient:       muxClient,
		options:         options,
		clientId:        0,
		pendingRequests: make(map[string]*PendingRequest),
		logger:          logger,
		serverInfo:      mcp.ServerInfo{},
		tools:           []mcp.ToolDescription{},
	}
}

func (c *MCPProxyClient) Start(ctx context.Context) error {
	transport := c.transport

	// TODO: this is not working as expected if we start
	// the proxy client before the mux server
	//
	// wait for the mux client to be started with a 10 second timeout
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	c.logger.Info("connecting to mux server...")

	for !c.muxClient.IsStarted() {
		select {
		case <-ticker.C:
			continue
		case <-timeout:
			return fmt.Errorf("timeout waiting for mux client to start")
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	c.logger.Info("connected to mux server")

	transport.OnMessage(func(msg json.RawMessage) {
		// check the message nature
		c.logger.Debug(fmt.Sprintf("received message: %s\n", string(msg)))
		nature, jsonRpcRawMessage, err := jsonrpc.CheckJsonMessage(msg)
		if err != nil {
			c.logger.Error(fmt.Sprintf("invalid received message: %s\n", string(msg)))
			return
		}

		// the MCP protocol does not support batch requests
		if nature == jsonrpc.MessageNatureBatchRequest {
			c.sendError(jsonrpc.RpcParseError, "response not supported yet", nil)
			return
		}
		// we process the message here
		c.handleIncomingMessage(jsonRpcRawMessage, nature)

	})

	// Set up error handler
	transport.OnError(func(err error) {
		fmt.Printf("[proxy] transport error: %s\n", err)
	})

	// Start the transport
	if err := transport.Start(ctx); err != nil {
		fmt.Printf("[proxy] failed to start transport: %s\n", err)
		return err
	}

	// First message to send is always an initialize request
	req, err := mkRpcRequestInitialize(c.options.ProxyName, version.Version, c.clientId)
	if err != nil {
		fmt.Printf("[proxy] failed to create initialize request: %s\n", err)
		return err
	}
	c.sendJsonRpcRequest(req)

	// Keep the main thread alive
	// will be interrupted by the context
	<-ctx.Done()

	transport.Close()

	c.logger.Info("shutdown\n")

	return nil
}

func (c *MCPProxyClient) sendJsonRpcRequest(request *jsonrpc.JsonRpcRequest) {
	jsonRequest, err := jsonrpc.MarshalJsonRpcRequest(request)
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to marshal request: %s\n", err))
		return
	}

	c.logger.Debug(fmt.Sprintf("sending request: %s\n", string(jsonRequest)))

	// send the message
	c.transport.Send(jsonRequest)

	// we store the request in the pending requests map
	// so we can match the response with the request
	c.pendingRequests[jsonrpc.RequestIdToString(request.Id)] = &PendingRequest{
		method:    request.Method,
		messageId: request.Id,
	}
	// increment the client id for the next message to send
	c.clientId++
}

func (c *MCPProxyClient) sendError(code int, message string, id *jsonrpc.JsonRpcRequestId) error {
	response := &jsonrpc.JsonRpcResponse{
		Error: &jsonrpc.JsonRpcError{
			Code:    code,
			Message: message,
		},
		Id: id,
	}
	jsonError, err := jsonrpc.MarshalJsonRpcResponse(response)
	if err != nil {
		return err
	}

	// send the message
	c.transport.Send(jsonError)

	return nil
}

func (c *MCPProxyClient) getPendingRequest(reqId *jsonrpc.JsonRpcRequestId) *PendingRequest {
	requestIdString := jsonrpc.RequestIdToString(reqId)
	return c.pendingRequests[requestIdString]
}
