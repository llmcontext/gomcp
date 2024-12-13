package jsonrpc

func NewJsonRpcNotification(method string) *JsonRpcRequest {
	return &JsonRpcRequest{
		JsonRpcVersion: JsonRpcVersion,
		Method:         method,
	}
}
