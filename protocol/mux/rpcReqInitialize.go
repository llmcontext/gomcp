package mux

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
)

const (
	RpcRequestMethodMuxInitialize = "mux/initialize"
)

type JsonRpcRequestMuxInitializeParams struct {
	ProtocolVersion string `json:"protocolVersion"`
	SessionId       string `json:"sessionId"`
}

func ParseJsonRpcRequestMuxInitializeParams(request *jsonrpc.JsonRpcRequest) (*JsonRpcRequestMuxInitializeParams, error) {
	// parse params
	if request.Params == nil {
		return nil, fmt.Errorf("missing params")
	}
	if !request.Params.IsNamed() {
		return nil, fmt.Errorf("params must be an object")
	}
	namedParams := request.Params.NamedParams

	req := JsonRpcRequestMuxInitializeParams{}

	// read protocol version
	protocolVersion, ok := namedParams["protocolVersion"].(string)
	if !ok {
		return nil, fmt.Errorf("missing protocolVersion")
	}
	req.ProtocolVersion = protocolVersion

	// read session id
	sessionId, ok := namedParams["sessionId"].(string)
	if !ok {
		return nil, fmt.Errorf("missing sessionId")
	}
	req.SessionId = sessionId

	return &req, nil
}
