package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/types"
)

const (
	GomcpProxyClientName    = "gomcp-proxy"
	GomcpProxyClientVersion = "0.1.0"
)

type PendingRequest struct {
	method    string
	messageId *jsonrpc.JsonRpcRequestId
}

type MCPProxyClient struct {
	transport types.Transport
	clientId  int
	logger    *ProxyLogger
	// pendingRequests is a map of message id to pending request
	pendingRequests map[string]*PendingRequest
}

func NewMCPProxyClient(transport types.Transport) *MCPProxyClient {
	return &MCPProxyClient{
		transport:       transport,
		clientId:        0,
		pendingRequests: make(map[string]*PendingRequest),
		logger:          NewProxyLogger(),
	}
}

func (c *MCPProxyClient) Start(ctx context.Context) error {
	transport := c.transport

	transport.OnMessage(func(msg json.RawMessage) {
		// check the message nature
		c.logger.Debug(fmt.Sprintf("received message: %s\n", string(msg)))
		nature, jsonRpcRawMessage, err := jsonrpc.CheckJsonMessage(msg)
		if err != nil {
			c.logger.Error(fmt.Sprintf("invalid received message: %s\n", string(msg)))
			os.Exit(1)
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
	req, err := mkRpcCallInitialize(GomcpProxyClientName, GomcpProxyClientVersion, c.clientId)
	if err != nil {
		fmt.Printf("[proxy] failed to create initialize request: %s\n", err)
		return err
	}
	c.sendJsonRpcRequest(req)

	// Keep the main thread alive
	// will be interrupted by the context
	<-ctx.Done()

	transport.Close()

	fmt.Printf("@@ [proxy] shutdown\n")

	return nil
}

func (c *MCPProxyClient) sendJsonRpcRequest(request *jsonrpc.JsonRpcRequest) {

	jsonRequest, err := jsonrpc.MarshalJsonRpcRequest(request)
	if err != nil {
		fmt.Printf("[proxy] failed to marshal request: %s\n", err)
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
