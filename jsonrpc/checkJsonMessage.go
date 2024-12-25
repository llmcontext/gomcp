package jsonrpc

import (
	"encoding/json"
	"fmt"
)

// detect if the message is either:
// - a notification
// - a request
// - a response
// - an error
// - a batch request
type MessageNature int

const (
	MessageNatureRequest MessageNature = iota
	MessageNatureResponse
	MessageNatureNotification
	MessageNatureBatchRequest
	MessageNatureUnknown
)

// specification at: https://www.jsonrpc.org/specification
// a request Object has:
// - jsonrpc: "2.0"
// - method: string
// - params: optional, object or array
// - id: number, string, null
//
// a notification has:
// - jsonrpc: "2.0"
// - method: string
// - params: optional, object or array
// (NO id field)
//
// a response has:
// - jsonrpc: "2.0"
// - result: optional, object or array
// - error: optional, object
//   - code: number
//   - data: optional, object or array
//   - message: string
// - id: number, string, null
//
// a batch request has:
// - array of request objects

// this function checks if the message is a valid JsonRpc message
// and returns the message nature
// it returns an error if the message is not a valid JsonRpc message
// this function is NOT a full validator, it just checks the general structure
// of the message to determine its nature
func CheckJsonMessage(message json.RawMessage) (MessageNature, JsonRpcRawMessage, error) {
	// we check if data is either an array or an object
	var rawJson interface{}
	if err := json.Unmarshal(message, &rawJson); err != nil {
		return MessageNatureUnknown, nil, err
	}

	switch v := rawJson.(type) {
	case []interface{}:
		return MessageNatureBatchRequest, nil, nil
	case map[string]interface{}:
		// any JsonRpc message must have a "jsonrpc" field
		if _, ok := v["jsonrpc"]; ok {
			if v["jsonrpc"] != JsonRpcVersion {
				return MessageNatureUnknown, nil, fmt.Errorf("invalid jsonrpc version: %s", v["jsonrpc"])
			}
		} else {
			return MessageNatureUnknown, nil, fmt.Errorf("missing jsonrpc field")
		}
		// any JsonRpc message must have a "method" field
		if _, ok := v["method"]; ok {
			// we check if the message is a request or a notification
			if _, ok := v["id"]; ok {
				return MessageNatureRequest, v, nil
			} else {
				return MessageNatureNotification, v, nil
			}
		} else {
			return MessageNatureResponse, v, nil
		}
	default:
		return MessageNatureUnknown, nil, fmt.Errorf("unknown message nature")
	}
}
