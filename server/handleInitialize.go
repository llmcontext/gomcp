package server

import (
	"encoding/json"

	"github.com/llmcontext/gomcp/jsonrpc"
)

type serverInitializeResponse struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    serverCapabilities `json:"capabilities"`
	ServerInfo      serverInfo         `json:"serverInfo"`
}

type serverInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// for now, only tools are supported
type serverCapabilities struct {
	Tools   serverCapabilitiesTools   `json:"tools"`
	Prompts serverCapabilitiesPrompts `json:"prompts"`
}

type serverCapabilitiesTools struct {
	ListChanged bool `json:"listChanged"`
}

type serverCapabilitiesPrompts struct {
	ListChanged bool `json:"listChanged"`
}

func (s *MCPServer) handleInitialize(request *jsonrpc.JsonRpcRequest) error {
	// parse params
	if request.Params == nil {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidRequest,
			Message: "missing params",
		}, request.Id)
	}
	if !request.Params.IsNamed() {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidRequest,
			Message: "params must be an object",
		}, request.Id)
	}
	namedParams := request.Params.NamedParams

	s.clientInfo = &ClientInfo{}

	protocolVersion, ok := namedParams["protocolVersion"].(string)
	if !ok {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidRequest,
			Message: "missing protocolVersion",
		}, request.Id)
	}
	s.protocolVersion = protocolVersion

	// read client information
	clientInfo, ok := namedParams["clientInfo"].(map[string]interface{})
	if !ok {
		s.sendError(&jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInvalidRequest,
			Message: "missing clientInfo",
		}, request.Id)
	}

	s.clientInfo.name = clientInfo["name"].(string)
	s.clientInfo.version = clientInfo["version"].(string)

	// prepare response
	response := serverInitializeResponse{
		ProtocolVersion: s.protocolVersion,
		Capabilities: serverCapabilities{
			Tools: serverCapabilitiesTools{
				ListChanged: false,
			},
			Prompts: serverCapabilitiesPrompts{
				ListChanged: false,
			},
		},
		ServerInfo: serverInfo{Name: s.serverName, Version: s.serverVersion},
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
