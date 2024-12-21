package jsonrpc

import (
	"fmt"
	"strconv"
	"strings"
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

func ReqIdStringToId(reqId string) *JsonRpcRequestId {
	if strings.HasPrefix(reqId, "N:") {
		number := strings.TrimPrefix(reqId, "N:")
		intValue, err := strconv.Atoi(number)
		if err != nil {
			return nil
		}
		return &JsonRpcRequestId{Number: &intValue}
	}
	if strings.HasPrefix(reqId, "S:") {
		stringValue := strings.TrimPrefix(reqId, "S:")
		return &JsonRpcRequestId{String: &stringValue}
	}
	return nil
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

type reqIdMappingEntry struct {
	reqIdA *JsonRpcRequestId
	reqIdB *JsonRpcRequestId
}

type ReqIdMapping struct {
	entries []reqIdMappingEntry
}

func NewReqIdMapping() *ReqIdMapping {
	return &ReqIdMapping{
		entries: []reqIdMappingEntry{},
	}
}

func (m *ReqIdMapping) AddMapping(reqIdA *JsonRpcRequestId, reqIdB *JsonRpcRequestId) {
	m.entries = append(m.entries, reqIdMappingEntry{reqIdA, reqIdB})
}

func (m *ReqIdMapping) GetMapping(reqId *JsonRpcRequestId) *JsonRpcRequestId {
	var foundEntry *JsonRpcRequestId = nil
	var ixEntry int = -1
	for ix, entry := range m.entries {
		// Compare the actual values instead of pointer addresses
		if equalRequestIds(entry.reqIdA, reqId) {
			foundEntry = entry.reqIdB
			ixEntry = ix
			break
		}
	}
	if foundEntry != nil {
		// delete the entry from the list
		m.entries = append(m.entries[:ixEntry], m.entries[ixEntry+1:]...)
	}
	return foundEntry
}

func equalRequestIds(a, b *JsonRpcRequestId) bool {
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
