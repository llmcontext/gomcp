package mcpClient

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
)

func mkRpcRequestInitialize(clientName string, clientVersion string, id int) (*jsonrpc.JsonRpcRequest, error) {
	// we create the parameters for the initialize request
	// the proxy does not have any capabilities
	params := mcp.JsonRpcRequestInitializeParams{
		ProtocolVersion: mcp.ProtocolVersion,
		Capabilities:    mcp.ClientCapabilities{},
		ClientInfo: mcp.ClientInfo{
			Name:    clientName,
			Version: clientVersion,
		},
	}

	// we create the JSON-RPC request
	req := jsonrpc.NewJsonRpcRequestWithNamedParams(
		mcp.RpcRequestMethodInitialize, params, id)

	if req == nil {
		return nil, fmt.Errorf("failed to create initialize request")
	}

	return req, nil
}

func mkRpcRequestToolsList(id int) (*jsonrpc.JsonRpcRequest, error) {
	params := mcp.JsonRpcRequestToolsListParams{}

	req := jsonrpc.NewJsonRpcRequestWithNamedParams(
		mcp.RpcRequestMethodToolsList, params, id)

	if req == nil {
		return nil, fmt.Errorf("failed to create tools list request")
	}

	return req, nil
}
