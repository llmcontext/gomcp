package jsonrpc

import "encoding/json"

const (
	JsonRpcVersion = "2.0"
)

type JsonRpcError struct {
	Code    int
	Message string
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
