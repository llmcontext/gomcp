package jsonrpc

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const debugParseRequest = false

// JSON-RPC 2.0 Specification
// https://www.jsonrpc.org/specification

type RawJsonRpcRequestMessage struct {
	JsonRpcVersion string      `json:"jsonrpc"`
	Method         string      `json:"method"`
	Params         interface{} `json:"params,omitempty"`
	Id             interface{} `json:"id,omitempty"`
}

type JsonRequestParseResponse struct {
	Request   *JsonRpcRequest
	RequestId *JsonRpcRequestId
	Error     *JsonRpcError
}

func ParseJsonRpcRequest(rawJson JsonRpcRawMessage) (*JsonRpcRequest, *JsonRpcRequestId, *JsonRpcError) {
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

func MarshalJsonRpcRequest(request *JsonRpcRequest) ([]byte, error) {
	var rawParams interface{} = nil
	if request.Params != nil {
		if request.Params.IsPositional() {
			rawParams = request.Params.PositionalParams
		} else {
			rawParams = request.Params.NamedParams
		}
	}

	var rawId interface{} = nil
	if request.Id != nil {
		if request.Id.Number != nil {
			rawId = *request.Id.Number
		} else if request.Id.String != nil {
			rawId = *request.Id.String
		}
	}

	rawJson := RawJsonRpcRequestMessage{
		JsonRpcVersion: JsonRpcVersion,
		Method:         request.Method,
		Params:         rawParams,
		Id:             rawId,
	}
	return json.Marshal(rawJson)
}
