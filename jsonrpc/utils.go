package jsonrpc

import (
	"encoding/json"
	"fmt"
)

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

func RequestIdToString(register *JsonRpcRequestId) string {
	if register == nil {
		return "**no-id**"
	}

	if register.Number != nil {
		return fmt.Sprintf("%d", *register.Number)
	}
	return *register.String
}

func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

// func rawMessagePtr(s string) *json.RawMessage {
// 	rm := json.RawMessage(s)
// 	return &rm
// }

func structToMap(obj interface{}) (map[string]interface{}, error) {
	// First marshal the struct to JSON
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	// Then unmarshal into a map
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return result, nil
}
