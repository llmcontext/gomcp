package mux

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
)

type JsonRpcResponseProxyRegisterResult struct {
	SessionId  string `json:"sessionId"`
	ProxyId    string `json:"proxyId"`
	Persistent bool   `json:"persistent"`
	Denied     bool   `json:"denied"`
}

func ParseJsonRpcResponseProxyRegister(response *jsonrpc.JsonRpcResponse) (*JsonRpcResponseProxyRegisterResult, error) {
	if response.Result == nil {
		return nil, fmt.Errorf("missing result")
	}
	// parse params
	result, err := protocol.CheckIsObject(response.Result, "result")
	if err != nil {
		return nil, err
	}

	// read sessionId
	sessionId, err := protocol.GetStringField(result, "sessionId")
	if err != nil {
		return nil, err
	}

	// read proxyId
	proxyId, err := protocol.GetStringField(result, "proxyId")
	if err != nil {
		return nil, err
	}

	// read persistent
	persistent, err := protocol.GetBoolField(result, "persistent")
	if err != nil {
		return nil, err
	}

	// read denied
	denied, err := protocol.GetBoolField(result, "denied")
	if err != nil {
		return nil, err
	}

	return &JsonRpcResponseProxyRegisterResult{
		SessionId:  sessionId,
		ProxyId:    proxyId,
		Persistent: persistent,
		Denied:     denied,
	}, nil
}
