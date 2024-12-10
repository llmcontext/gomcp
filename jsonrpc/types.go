package jsonrpc

import "encoding/json"

const (
	JsonRpcVersion = "2.0"
)

const (
	// JSON-RPC 2.0 Specification
	// https://www.jsonrpc.org/specification#error_object
	RpcParseError     = -32700
	RpcInvalidRequest = -32600
	RpcMethodNotFound = -32601
	RpcInvalidParams  = -32602
	RpcInternalError  = -32603
)

// most generic type for a JsonRpc message that is not a batch request
type JsonRpcRawMessage map[string]interface{}

type JsonRpcError struct {
	Code    int
	Message string
	Data    *json.RawMessage
}

type JsonRpcRequestId struct {
	Number *int
	String *string
}

type JsonRpcParams struct {
	PositionalParams []interface{}
	NamedParams      map[string]interface{}
}

func (p *JsonRpcParams) IsPositional() bool {
	return p.PositionalParams != nil
}

func (p *JsonRpcParams) IsNamed() bool {
	return p.NamedParams != nil
}

type JsonRpcRequest struct {
	JsonRpcVersion string
	Method         string
	Params         *JsonRpcParams
	Id             *JsonRpcRequestId
}

type JsonRpcResponse struct {
	Result *json.RawMessage
	Error  *JsonRpcError
	Id     *JsonRpcRequestId
}

func NewJsonRpcRequestWithNamedParams(method string, params interface{}, id int) *JsonRpcRequest {
	// First marshal to JSON bytes
	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil
	}

	// Then unmarshal into map[string]interface{}
	var namedParams map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &namedParams); err != nil {
		return nil
	}

	return &JsonRpcRequest{
		JsonRpcVersion: JsonRpcVersion,
		Method:         method,
		Params:         &JsonRpcParams{NamedParams: namedParams},
		Id:             &JsonRpcRequestId{Number: &id},
	}
}
