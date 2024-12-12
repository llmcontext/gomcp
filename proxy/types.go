package proxy

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/jsonrpc/messages"
)

const (
	ProtocolVersion = "2024-11-05"
	JsonRpcVersion  = "2.0"
)

func mkRpcRequestInitialize(clientName string, clientVersion string, id int) (*jsonrpc.JsonRpcRequest, error) {
	// we create the parameters for the initialize request
	// the proxy does not have any capabilities
	params := messages.JsonRpcRequestInitialize{
		ProtocolVersion: ProtocolVersion,
		Capabilities:    messages.ClientCapabilities{},
		ClientInfo: messages.ClientInfo{
			Name:    clientName,
			Version: clientVersion,
		},
	}

	// we create the JSON-RPC request
	req := jsonrpc.NewJsonRpcRequestWithNamedParams(
		messages.RpcRequestMethodInitialize, params, id)

	if req == nil {
		return nil, fmt.Errorf("failed to create initialize request")
	}

	return req, nil
}

func mkRpcNotification(method string) (*jsonrpc.JsonRpcRequest, error) {
	// we create the JSON-RPC request
	req := jsonrpc.NewJsonRpcNotification(method)

	if req == nil {
		return nil, fmt.Errorf("failed to create notification")
	}

	return req, nil
}
