package mcp

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
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
	protocolVersion, ok := namedParams["protocolVersion"].(string)
	if !ok {
		return nil, fmt.Errorf("missing protocolVersion")
	}
	req.ProtocolVersion = protocolVersion

	// read client information
	clientInfo, ok := namedParams["clientInfo"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing clientInfo")
	}
	name, ok := clientInfo["name"].(string)
	if !ok {
		return nil, fmt.Errorf("clientInfo.name must be a string")
	}
	req.ClientInfo.Name = name
	version, ok := clientInfo["version"].(string)
	if !ok {
		return nil, fmt.Errorf("clientInfo.version must be a string")
	}
	req.ClientInfo.Version = version

	// TODO: read capabilities and roots

	return &req, nil
}
