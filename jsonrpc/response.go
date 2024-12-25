package jsonrpc

import (
	"encoding/json"
	"errors"
	"fmt"
)

type RawJsonRpcResponseMessage struct {
	JsonRpcVersion string           `json:"jsonrpc"`
	Result         *json.RawMessage `json:"result,omitempty"`
	Id             json.RawMessage  `json:"id"`
}

type RawJsonError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RawJsonRpcErrorMessage struct {
	JsonRpcVersion string          `json:"jsonrpc"`
	Error          RawJsonError    `json:"error"`
	Id             json.RawMessage `json:"id"`
}

func parseRpcId(id *JsonRpcRequestId) json.RawMessage {
	var responseId json.RawMessage
	if id != nil {
		if id.String != nil {
			responseId, _ = json.Marshal(*id.String)
		} else if id.Number != nil {
			responseId, _ = json.Marshal(*id.Number)
		} else {
			responseId = nil
		}
	} else {
		responseId = nil
	}
	return responseId
}

func MarshalJsonRpcResponse(response *JsonRpcResponse) ([]byte, error) {
	if response.Error == nil && response.Result == nil {
		return nil, errors.New("error or result is required")
	}

	if response.Error != nil && response.Result != nil {
		return nil, errors.New("error and result cannot be set at the same time")
	}

	// we marshal the id to a json.RawMessage
	responseId := parseRpcId(response.Id)

	// check if we have an error message
	if response.Error != nil {
		return json.Marshal(RawJsonRpcErrorMessage{
			JsonRpcVersion: JsonRpcVersion,
			Error:          RawJsonError{Code: response.Error.Code, Message: response.Error.Message},
			Id:             responseId,
		})
	}

	// we marshal the result to a json.RawMessage
	result, err := json.Marshal(response.Result)
	if err != nil {
		return nil, err
	}
	rawResult := json.RawMessage(result)

	// we have a result message
	return json.Marshal(RawJsonRpcResponseMessage{
		JsonRpcVersion: JsonRpcVersion,
		Result:         &rawResult,
		Id:             responseId,
	})
}

func ParseJsonRpcResponse(rawMessage JsonRpcRawMessage) (*JsonRpcResponse, *JsonRpcRequestId, *JsonRpcError) {
	// we initialize the request with default values
	requestId := extractId(rawMessage)
	jsonRpcResponse := &JsonRpcResponse{
		JsonRpcVersion: "",
		Result:         nil,
		Error:          nil,
		Id:             requestId,
	}

	// we iterate through all the keys in the JSON object
	// and check if they are valid
	for key, value := range rawMessage {
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
			jsonRpcResponse.JsonRpcVersion = version

		case "result":
			jsonRpcResponse.Result = value

		case "error":
			// we check that this is a map[string]interface{}
			error, ok := value.(map[string]interface{})
			if !ok {
				return nil, requestId, &JsonRpcError{
					Code:    RpcParseError,
					Message: "invalid error",
				}
			}
			jsonRpcResponse.Error = parseJsonRpcError(error)

		case "id":
			// already parsed

		default:
			return nil, requestId, &JsonRpcError{
				Code:    RpcInvalidRequest,
				Message: fmt.Sprintf("invalid key: %s", key),
			}

		}
	}

	return jsonRpcResponse, requestId, nil

}

// no strict parsing of the error
func parseJsonRpcError(error map[string]interface{}) *JsonRpcError {
	jsonRpcError := &JsonRpcError{
		Code:    0,
		Message: "",
	}

	for key, value := range error {
		switch key {
		case "code":
			code, ok := value.(float64)
			if !ok {
				return nil
			}
			jsonRpcError.Code = int(code)
		case "message":
			message, ok := value.(string)
			if !ok {
				return nil
			}
			jsonRpcError.Message = message
		}
	}
	return jsonRpcError
}
