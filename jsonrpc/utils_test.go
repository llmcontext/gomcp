package jsonrpc

import (
	"testing"
)

func TestReqIdMapping(t *testing.T) {
	tests := []struct {
		name     string
		reqIdA   *JsonRpcRequestId
		reqIdB   *JsonRpcRequestId
		wantNil  bool // for testing non-existent mappings
		twoReads bool // for testing removal after retrieval
	}{
		{
			name:   "number request IDs",
			reqIdA: &JsonRpcRequestId{Number: intPtr(1)},
			reqIdB: &JsonRpcRequestId{Number: intPtr(2)},
		},
		{
			name:   "string request IDs",
			reqIdA: &JsonRpcRequestId{String: stringPtr("test1")},
			reqIdB: &JsonRpcRequestId{String: stringPtr("test2")},
		},
		{
			name:     "remove mapping after retrieval",
			reqIdA:   &JsonRpcRequestId{Number: intPtr(1)},
			reqIdB:   &JsonRpcRequestId{Number: intPtr(2)},
			twoReads: true,
		},
		{
			name:   "nil request IDs",
			reqIdA: nil,
			reqIdB: &JsonRpcRequestId{Number: intPtr(1)},
		},
		{
			name:    "non-existent mapping",
			reqIdA:  &JsonRpcRequestId{Number: intPtr(1)},
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapping := NewReqIdMapping()

			if !tt.wantNil {
				mapping.AddMapping(tt.reqIdA, tt.reqIdB)
			}

			result := mapping.GetMapping(tt.reqIdA)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected nil for non-existent mapping, got %v", RequestIdToString(result))
				}
				return
			}

			if !equalRequestIds(result, tt.reqIdB) {
				t.Errorf("Expected request ID %v, got %v", RequestIdToString(tt.reqIdB), RequestIdToString(result))
			}

			if tt.twoReads {
				// Second retrieval should return nil
				result2 := mapping.GetMapping(tt.reqIdA)
				if result2 != nil {
					t.Errorf("Expected second retrieval to return nil, got %v", RequestIdToString(result2))
				}
			}
		})
	}
}
