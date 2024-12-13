package muxClient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type MuxClient struct {
	muxAddress string
	logger     types.TermLogger
	transport  types.Transport
	isStarted  bool
}

func NewMuxClient(muxAddress string, logger types.TermLogger) *MuxClient {
	return &MuxClient{
		muxAddress: muxAddress,
		logger:     logger,
		transport:  transport.NewSocketClientTransport(muxAddress),
		isStarted:  false,
	}
}

func (c *MuxClient) IsStarted() bool {
	return c.isStarted
}

func (c *MuxClient) Start(ctx context.Context) error {
	// Set up message handling
	c.transport.OnMessage(func(msg json.RawMessage) {
		// Handle incoming messages here
		// You might want to decode the message and process it accordingly
		fmt.Printf("Received message: %v\n", msg)
	})

	// Set up error handling
	c.transport.OnError(func(err error) {
		// Handle transport errors
		fmt.Printf("Transport error: %v\n", err)
	})

	// Set up close handling
	c.transport.OnClose(func() {
		// Handle transport closure
		fmt.Println("Transport closed")
	})

	// Start the transport with context
	if err := c.transport.Start(ctx); err != nil {
		return fmt.Errorf("failed to start transport: %w", err)
	}

	c.isStarted = true

	return nil
}

// SendMessage sends a JSON-encodable message through the transport
func (c *MuxClient) SendMessage(ctx context.Context, message interface{}) error {
	// Convert message to JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Send the message
	if err := c.transport.Send(jsonData); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (c *MuxClient) GetTransport() types.Transport {
	return c.transport
}
