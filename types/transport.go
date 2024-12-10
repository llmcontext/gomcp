package types

import (
	"context"
	"encoding/json"
)

// Transport defines the interface for MCP communication
type Transport interface {
	// Starts processing messages on the transport, including any connection
	// steps that might need to be taken.
	//
	// This method should only be called after callbacks are installed, or else
	// messages may be lost.
	//
	// NOTE: This method should not be called explicitly when using Client,
	// Server, or Protocol classes, as they will implicitly call start().
	Start(ctx context.Context) error

	// Sends a JSON-RPC message (request or response).
	Send(message json.RawMessage) error

	// Callback for when a message (request or response) is received over the
	// connection.
	OnMessage(callback func(json.RawMessage))

	// Closes the connection.
	Close()

	// Callback for when the connection is closed for any reason.
	//
	// This should be invoked when close() is called as well.
	OnClose(callback func())

	// Callback for when an error occurs.
	//
	// Note that errors are not necessarily fatal; they are used for reporting
	// any kind of exceptional condition out of band.
	OnError(callback func(error))
}
