package jsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const debugParseRequest = false

// JSON-RPC 2.0 Specification
// https://www.jsonrpc.org/specification

const (
	// JSON-RPC 2.0 Specification
	// https://www.jsonrpc.org/specification#error_object
	RpcParseError     = -32700
	RpcInvalidRequest = -32600
	RpcMethodNotFound = -32601
	RpcInvalidParams  = -32602
	RpcInternalError  = -32603
)

type JsonRequestParseResponse struct {
	Request   *JsonRpcRequest
	RequestId *JsonRpcRequestId
	Error     *JsonRpcError
}

// returns a list of requests,
// a boolean indicating if the request is a batch, and an error
func ParseRequest(data []byte) ([]JsonRequestParseResponse, bool, *JsonRpcError) {
	// we check if data is either an array or an object
	var rawJson interface{}
	if err := json.Unmarshal(data, &rawJson); err != nil {
		return nil, false, &JsonRpcError{
			Code:    RpcParseError,
			Message: err.Error(),
		}
	}

	switch v := rawJson.(type) {
	case []interface{}:
		// we parse each element in the array as a request
		requests := []JsonRequestParseResponse{}
		for _, element := range v {
			if element, ok := element.(map[string]interface{}); ok {
				// we check if the element is a map[string]interface{}
				request, requestId, err := ParseSimpleRequest(element)
				requests = append(requests, JsonRequestParseResponse{Request: request, RequestId: requestId, Error: err})
			} else {
				return nil, false, &JsonRpcError{
					Code:    RpcInvalidRequest,
					Message: "request in batch must be an object",
				}
			}
		}
		return requests, true, nil

	case map[string]interface{}:
		// if it is an object, we parse it as a simple request
		request, requestId, err := ParseSimpleRequest(v)
		return []JsonRequestParseResponse{{Request: request, RequestId: requestId, Error: err}}, false, nil

	default:
		return nil, false, &JsonRpcError{
			Code:    RpcInvalidRequest,
			Message: "request must be an object or an array",
		}
	}
}

func ParseSimpleRequest(rawJson map[string]interface{}) (*JsonRpcRequest, *JsonRpcRequestId, *JsonRpcError) {
	// we initialize the request with default values
	requestId := extractId(rawJson)
	jsonRpcRequest := &JsonRpcRequest{
		JsonRpcVersion: "",
		Method:         "",
		Params:         nil,
		Id:             requestId,
	}

	// we iterate through all the keys in the JSON object
	// and check if they are valid
	for key, value := range rawJson {
		if debugParseRequest {
			fmt.Printf("key: %s, value: %v [type:%s]\n", key, value, reflect.TypeOf(value))
		}

		switch key {
		case "jsonrpc":
			// we assert that the value is a string
			version, ok := value.(string)
			if !ok {
				return nil, requestId, &JsonRpcError{
					Code:    RpcParseError,
					Message: "invalid JSON-RPC version",
				}
			}
			jsonRpcRequest.JsonRpcVersion = version

		case "method":
			if value == nil {
				return nil, requestId, &JsonRpcError{
					Code:    RpcInvalidRequest,
					Message: "method is required",
				}
			}
			// we assert that the value is a string
			method, ok := value.(string)
			if !ok {
				return nil, requestId, &JsonRpcError{
					Code:    RpcInvalidRequest,
					Message: "method must be a string",
				}
			}
			jsonRpcRequest.Method = method
		case "params":
			// the params can be a positional array or a named object
			switch v := value.(type) {
			case []interface{}:
				jsonRpcRequest.Params = &JsonRpcParams{PositionalParams: v, NamedParams: nil}
			case map[string]interface{}:
				jsonRpcRequest.Params = &JsonRpcParams{NamedParams: v, PositionalParams: nil}
			default:
				return nil, requestId, &JsonRpcError{
					Code:    RpcInvalidRequest,
					Message: "params must be an array or an object",
				}
			}
		case "id":
			// already parsed

		default:
			return nil, requestId, &JsonRpcError{
				Code:    RpcInvalidRequest,
				Message: fmt.Sprintf("invalid key: %s", key),
			}
		}
	}

	// check if the jsonrpc version is empty
	if jsonRpcRequest.JsonRpcVersion != JsonRpcVersion {
		return nil, requestId, &JsonRpcError{
			Code:    RpcInvalidRequest,
			Message: "invalid JSON-RPC version",
		}
	}

	// check if the method is empty
	if jsonRpcRequest.Method == "" {
		return nil, requestId, &JsonRpcError{
			Code:    RpcInvalidRequest,
			Message: "method is required",
		}
	}

	return jsonRpcRequest, requestId, nil

}

func extractId(rawJson map[string]interface{}) *JsonRpcRequestId {
	value, ok := rawJson["id"]
	if !ok {
		return nil
	}

	// the id can be a number or a string
	switch v := value.(type) {
	case int:
		return &JsonRpcRequestId{Number: &v}
	case float64: // Changed from int to float64
		intValue := int(v)
		return &JsonRpcRequestId{Number: &intValue}
	case string:
		return &JsonRpcRequestId{String: &v}
	default:
		return nil
	}
}
