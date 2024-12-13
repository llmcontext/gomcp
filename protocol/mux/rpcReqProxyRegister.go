package mux

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

const (
	RpcRequestMethodProxyRegister = "proxy/register"
)

type JsonRpcRequestProxyRegisterParams struct {
	ProtocolVersion string           `json:"protocolVersion"`
	Proxy           ProxyDescription `json:"proxy"`
	ServerInfo      ServerInfo       `json:"serverInfo"`
}

type ProxyDescription struct {
	WorkingDirectory string `json:"workingDirectory"`
	Command          string `json:"command"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func ParseJsonRpcRequestProxyRegisterParams(request *jsonrpc.JsonRpcRequest) (*JsonRpcRequestProxyRegisterParams, error) {
	// parse params
	if request.Params == nil {
		return nil, fmt.Errorf("missing params")
	}
	if !request.Params.IsNamed() {
		return nil, fmt.Errorf("params must be an object")
	}
	namedParams := request.Params.NamedParams

	req := JsonRpcRequestProxyRegisterParams{}

	// read protocol version
	protocolVersion, ok := namedParams["protocolVersion"].(string)
	if !ok {
		return nil, fmt.Errorf("missing protocolVersion")
	}
	req.ProtocolVersion = protocolVersion

	// read proxy
	proxy, err := protocol.GetObjectField(namedParams, "proxy")
	if err != nil {
		return nil, fmt.Errorf("missing proxy")
	}
	req.Proxy.WorkingDirectory, err = protocol.GetStringField(proxy, "workingDirectory")
	if err != nil {
		return nil, fmt.Errorf("proxy.workingDirectory must be a string")
	}
	req.Proxy.Command, err = protocol.GetStringField(proxy, "command")
	if err != nil {
		return nil, fmt.Errorf("proxy.command must be a string")
	}

	return &req, nil
}
