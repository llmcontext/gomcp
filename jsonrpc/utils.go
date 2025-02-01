package jsonrpc

import (
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
		return "X"
	}

	if register.Number != nil {
		return fmt.Sprintf("N:%d", *register.Number)
	}
	return fmt.Sprintf("S:%s", *register.String)
}

func StringPtr(s string) *string {
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func BoolPtr(b bool) *bool {
	return &b
}

func EqualRequestIds(a, b *JsonRpcRequestId) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a.Number != nil && b.Number != nil {
		return *a.Number == *b.Number
	}
	if a.String != nil && b.String != nil {
		return *a.String == *b.String
	}
	return false
}
