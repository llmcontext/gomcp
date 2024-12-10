package proxy

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
)

const (
	ProtocolVersion = "2024-11-05"
	JsonRpcVersion  = "2.0"
)

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ParamsInitialize struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

func mkRpcCallInitialize(clientName string, clientVersion string, id int) (*jsonrpc.JsonRpcRequest, error) {
	params := ParamsInitialize{
		ProtocolVersion: ProtocolVersion,
		Capabilities:    map[string]interface{}{},
		ClientInfo: ClientInfo{
			Name:    clientName,
			Version: clientVersion,
		},
	}

	req := jsonrpc.NewJsonRpcRequestWithNamedParams("initialize", params, id)

	if req == nil {
		return nil, fmt.Errorf("failed to create initialize request")
	}

	return req, nil
}
