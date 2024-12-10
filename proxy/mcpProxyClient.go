package proxy

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/types"
)

const (
	ClientName    = "gomcp-proxy"
	ClientVersion = "0.1.0"
)

type PendingRequest struct {
	method    string
	messageId *jsonrpc.JsonRpcRequestId
}

type MCPProxyClient struct {
	transport types.Transport
	clientId  int
	// pendingRequests is a map of message id to pending request
	pendingRequests map[string]PendingRequest
}

func NewMCPProxyClient(transport types.Transport) *MCPProxyClient {
	return &MCPProxyClient{
		transport:       transport,
		clientId:        0,
		pendingRequests: make(map[string]PendingRequest),
	}
}

func (c *MCPProxyClient) Start(ctx context.Context) error {
	transport := c.transport

	transport.OnMessage(func(msg json.RawMessage) {
		// check the message nature
		nature, jsonRpcRawMessage, err := jsonrpc.CheckJsonMessage(msg)
		if err != nil {
			c.sendError(jsonrpc.RpcParseError, err.Error(), nil)
			fmt.Printf("@@ [proxy] invalid transport received message: %s\n", string(msg))
			return
		}

		if nature == jsonrpc.MessageNatureBatchRequest {
			c.sendError(jsonrpc.RpcParseError, "response not supported yet", nil)
			return
		}
		fmt.Printf("[proxy] received request: %+v\n", jsonRpcRawMessage)

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
	req, err := mkRpcCallInitialize(ClientName, ClientVersion, c.clientId)
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

	fmt.Printf("@@ [proxy] sending request: %s\n", string(jsonRequest))

	// send the message
	c.transport.Send(jsonRequest)

	c.pendingRequests[requestIdToString(request.Id)] = PendingRequest{
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

func requestIdToString(id *jsonrpc.JsonRpcRequestId) string {
	if id.Number != nil {
		return strconv.Itoa(*id.Number)
	} else if id.String != nil {
		return *id.String
	}
	return "*"
}
