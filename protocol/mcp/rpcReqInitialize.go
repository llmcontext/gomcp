package mcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

// specification
// https://spec.modelcontextprotocol.io/specification/basic/lifecycle/#initialization

const (
	RpcRequestMethodInitialize = "initialize"
)

type JsonRpcRequestInitializeParams struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo         `json:"clientInfo"`
}

type ClientCapabilities struct {
	Roots    ClientCapabilitiesRoots    `json:"roots"`
	Sampling ClientCapabilitiesSampling `json:"sampling"`
}

type ClientCapabilitiesRoots struct {
	ListChanged bool `json:"listChanged"`
}

type ClientCapabilitiesSampling struct {
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func ParseJsonRpcRequestInitialize(request *jsonrpc.JsonRpcRequest) (*JsonRpcRequestInitializeParams, error) {
	req := JsonRpcRequestInitializeParams{}

	// parse params
	if request.Params == nil {
		return nil, fmt.Errorf("missing params")
	}
	if !request.Params.IsNamed() {
		return nil, fmt.Errorf("params must be an object")
	}
	namedParams := request.Params.NamedParams

	// read protocol version
	protocolVersion, err := protocol.GetStringField(namedParams, "protocolVersion")
	if err != nil {
		return nil, fmt.Errorf("missing protocolVersion")
	}
	req.ProtocolVersion = protocolVersion

	// read client information
	clientInfo, err := protocol.GetObjectField(namedParams, "clientInfo")
	if err != nil {
		return nil, fmt.Errorf("missing clientInfo")
	}
	name, err := protocol.GetStringField(clientInfo, "name")
	if err != nil {
		return nil, fmt.Errorf("clientInfo.name must be a string")
	}
	req.ClientInfo.Name = name
	version, err := protocol.GetStringField(clientInfo, "version")
	if err != nil {
		return nil, fmt.Errorf("clientInfo.version must be a string")
	}
	req.ClientInfo.Version = version

	// TODO: read capabilities and roots

	return &req, nil
}
