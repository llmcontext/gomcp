package proxy

import (
	"encoding/json"
	"fmt"
)

const (
	ProtocolVersion = "2024-11-05"
	JsonRpcVersion  = "2.0"
)

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type JsonRpcCall struct {
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
	JsonRpc string      `json:"jsonrpc"`
	Id      int         `json:"id"`
}

type ParamsInitialize struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

func mkRpcCallInitialize(clientName string, clientVersion string, id int) json.RawMessage {
	params := ParamsInitialize{
		ProtocolVersion: ProtocolVersion,
		Capabilities:    map[string]interface{}{},
		ClientInfo: ClientInfo{
			Name:    clientName,
			Version: clientVersion,
		},
	}

	req := JsonRpcCall{
		JsonRpc: JsonRpcVersion,
		Method:  "initialize",
		Params:  params,
		Id:      id,
	}

	// marshal to json
	json, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("[proxy] failed to marshal initialize request: %s\n", err)
		return nil
	}
	return json
}
