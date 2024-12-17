package hubmcpserver

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
)

func (s *MCPServer) handleInitialize(request *jsonrpc.JsonRpcRequest) error {
	// parse params
	parsed, err := mcp.ParseJsonRpcRequestInitialize(request)
	if err != nil {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidRequest,
			Message: err.Error(),
		}, request.Id)
	}

	// store client information
	s.protocolVersion = parsed.ProtocolVersion
	s.clientInfo = &ClientInfo{
		name:    parsed.ClientInfo.Name,
		version: parsed.ClientInfo.Version,
	}

	// prepare response
	response := mcp.JsonRpcResponseInitializeResult{
		ProtocolVersion: s.protocolVersion,
		Capabilities: mcp.ServerCapabilities{
			Tools: &mcp.ServerCapabilitiesTools{
				ListChanged: jsonrpc.BoolPtr(false),
			},
			Prompts: &mcp.ServerCapabilitiesPrompts{
				ListChanged: jsonrpc.BoolPtr(false),
			},
		},
		ServerInfo: mcp.ServerInfo{Name: s.serverName, Version: s.serverVersion},
	}

	// marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInternalError,
			Message: "failed to marshal response",
		}, request.Id)
	}
	jsonResponse := json.RawMessage(responseBytes)

	// send response
	s.sendResponse(&jsonrpc.JsonRpcResponse{
		Id:     request.Id,
		Result: &jsonResponse,
	})

	return nil
}
