package jsonrpc

import (
	"encoding/json"
	"errors"
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

	// we have a result message
	return json.Marshal(RawJsonRpcResponseMessage{
		JsonRpcVersion: JsonRpcVersion,
		Result:         response.Result,
		Id:             responseId,
	})
}
