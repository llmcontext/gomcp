package proxy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/types"
)

type MCPProxyClient struct {
	transport types.Transport
}

func NewMCPProxyClient(transport types.Transport) *MCPProxyClient {
	return &MCPProxyClient{
		transport: transport,
	}
}

func (c *MCPProxyClient) Start(ctx context.Context) error {
	transport := c.transport

	transport.OnMessage(func(msg json.RawMessage) {
		requests, isBatch, error := jsonrpc.ParseRequest(msg)
		if error != nil {
			fmt.Printf("@@ [proxy] invalid transport received message: %s\n", string(msg))
			return
		}

		if isBatch {
			fmt.Printf("@@ [proxy] batch request not supported yet\n")
			return
		}

		request := requests[0]
		if request.Error != nil {
			fmt.Printf("@@ [proxy] error: %v\n", request.Error)
			return
		}

		fmt.Printf("@@ [proxy] received request: %+v\n", request)
	})

	// Set up error handler
	transport.OnError(func(err error) {
		fmt.Printf("@@ [proxy] transport error: %s\n", err)
	})

	// Start the transport
	if err := transport.Start(ctx); err != nil {
		fmt.Printf("@@ [proxy] failed to start transport: %s\n", err)
		return err
	}

	transport.Send(json.RawMessage(`{"jsonrpc": "2.0", "method": "ping", "id": 1}`))

	// Keep the main thread alive
	// will be interrupted by the context
	<-ctx.Done()

	transport.Close()

	fmt.Printf("@@ [proxy] shutdown\n")

	return nil
}
