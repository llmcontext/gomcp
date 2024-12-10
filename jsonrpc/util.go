package jsonrpc

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
