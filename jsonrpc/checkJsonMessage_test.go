package jsonrpc_test

import (
	"encoding/json"
	"testing"

	"github.com/llmcontext/gomcp/jsonrpc"
)

func TestCheckJsonMessage(t *testing.T) {
	tests := []struct {
		name       string
		message    string
		wantNature jsonrpc.MessageNature
		wantErr    bool
	}{
		{
			name:       "valid message request 1",
			message:    `{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			wantNature: jsonrpc.MessageNatureRequest,
			wantErr:    false,
		},
		{
			name:       "valid message request 2",
			message:    `{"jsonrpc": "2.0", "method": "subtract", "params": {"subtrahend": 23, "minuend": 42}, "id": 3}`,
			wantNature: jsonrpc.MessageNatureRequest,
			wantErr:    false,
		},
		{
			name:       "valid message notification",
			message:    `{"jsonrpc": "2.0", "method": "update", "params": [1,2,3,4,5]}`,
			wantNature: jsonrpc.MessageNatureNotification,
			wantErr:    false,
		},
		{
			name:       "valid message batch request",
			message:    `[{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}, {"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 2}]`,
			wantNature: jsonrpc.MessageNatureBatchRequest,
			wantErr:    false,
		},
		{
			name:       "valid response",
			message:    `{"jsonrpc": "2.0", "result": 19, "id": 1}`,
			wantNature: jsonrpc.MessageNatureResponse,
			wantErr:    false,
		},
		{
			name:       "valid response with error",
			message:    `{"jsonrpc": "2.0", "error": {"code": -32601, "message": "Method not found"}, "id": "1"}`,
			wantNature: jsonrpc.MessageNatureResponse,
			wantErr:    false,
		},
		{
			name:       "invalid JSON",
			message:    `{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1`,
			wantNature: jsonrpc.MessageNatureUnknown,
			wantErr:    true,
		},
		{
			name:       "missing jsonrpc field",
			message:    `{"method": "subtract", "params": [42, 23], "id": 1}`,
			wantNature: jsonrpc.MessageNatureUnknown,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nature, _, err := jsonrpc.CheckJsonMessage(json.RawMessage(tt.message))
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckJsonMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if nature != tt.wantNature {
				t.Errorf("CheckJsonMessage() nature = %v, wantNature %v", nature, tt.wantNature)
			}
		})
	}
}
