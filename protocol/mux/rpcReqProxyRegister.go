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
	ProtocolVersion string            `json:"protocolVersion"`
	Proxy           ProxyDescription  `json:"proxy"`
	ServerInfo      ServerInfo        `json:"serverInfo"`
	Tools           []ToolDescription `json:"tools"`
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

type ToolDescription struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
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

	req := JsonRpcRequestProxyRegisterParams{
		Tools: []ToolDescription{},
	}

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

	// read tools
	tools, err := protocol.GetArrayField(namedParams, "tools")
	if err != nil {
		return nil, fmt.Errorf("missing tools")
	}
	for _, tool := range tools {
		toolMap, ok := tool.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("tool must be an object")
		}
		toolDescription := ToolDescription{}
		toolDescription.Name, err = protocol.GetStringField(toolMap, "name")
		if err != nil {
			return nil, fmt.Errorf("tool.name must be a string")
		}
		toolDescription.Description, err = protocol.GetStringField(toolMap, "description")
		if err != nil {
			return nil, fmt.Errorf("tool.description must be a string")
		}
		toolDescription.InputSchema, err = protocol.GetObjectField(toolMap, "inputSchema")
		if err != nil {
			return nil, fmt.Errorf("tool.inputSchema must be an object")
		}
		req.Tools = append(req.Tools, toolDescription)
	}

	return &req, nil
}
