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
	ProxyId         string           `json:"proxyId"`
	Persistent      bool             `json:"persistent"`
	Proxy           ProxyDescription `json:"proxy"`
	ServerInfo      ServerInfo       `json:"serverInfo"`
}

type ProxyDescription struct {
	WorkingDirectory string   `json:"workingDirectory"`
	Command          string   `json:"command"`
	Args             []string `json:"args"`
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
	protocolVersion, err := protocol.GetStringField(namedParams, "protocolVersion")
	if err != nil {
		return nil, fmt.Errorf("missing protocolVersion")
	}
	req.ProtocolVersion = protocolVersion

	// read proxy id
	proxyId, err := protocol.GetStringField(namedParams, "proxyId")
	if err != nil {
		return nil, fmt.Errorf("missing proxyId")
	}
	req.ProxyId = proxyId

	// read persistent
	persistent, err := protocol.GetBoolField(namedParams, "persistent")
	if err != nil {
		return nil, fmt.Errorf("missing persistent")
	}
	req.Persistent = persistent

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

	// read server info
	serverInfo, err := protocol.GetObjectField(namedParams, "serverInfo")
	if err != nil {
		return nil, fmt.Errorf("missing serverInfo")
	}
	req.ServerInfo.Name, err = protocol.GetStringField(serverInfo, "name")
	if err != nil {
		return nil, fmt.Errorf("serverInfo.name must be a string")
	}
	req.ServerInfo.Version, err = protocol.GetStringField(serverInfo, "version")
	if err != nil {
		return nil, fmt.Errorf("serverInfo.version must be a string")
	}

	return &req, nil
}
