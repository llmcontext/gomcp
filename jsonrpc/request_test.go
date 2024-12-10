package jsonrpc

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseSimpleRequest(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantRequest *JsonRpcRequest
		wantError   *JsonRpcError
	}{
		{
			name:  "valid request with positional params",
			input: `{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			wantRequest: &JsonRpcRequest{
				JsonRpcVersion: "2.0",
				Method:         "subtract",
				Params:         &JsonRpcParams{PositionalParams: []interface{}{float64(42), float64(23)}},
				Id:             &JsonRpcRequestId{Number: intPtr(1)},
			},
			wantError: nil,
		},
		{
			name:  "valid request with named params",
			input: `{"jsonrpc": "2.0", "method": "subtract", "params": {"minuend": 42, "subtrahend": 23}, "id": "1"}`,
			wantRequest: &JsonRpcRequest{
				JsonRpcVersion: "2.0",
				Method:         "subtract",
				Params:         &JsonRpcParams{NamedParams: map[string]interface{}{"minuend": float64(42), "subtrahend": float64(23)}},
				Id:             &JsonRpcRequestId{String: stringPtr("1")},
			},
			wantError: nil,
		},
		{
			name:        "invalid JSON",
			input:       `{"jsonrpc": "2.0", "method": "subtract"`,
			wantRequest: nil,
			wantError: &JsonRpcError{
				Code:    RpcParseError,
				Message: "unexpected end of JSON input",
			},
		},
		{
			name:        "missing jsonrpc version",
			input:       `{"method": "subtract", "params": [42, 23], "id": 1}`,
			wantRequest: nil,
			wantError: &JsonRpcError{
				Code:    RpcInvalidRequest,
				Message: "invalid JSON-RPC version",
			},
		},
		{
			name:        "wrong jsonrpc version",
			input:       `{"jsonrpc": "1.0", "method": "subtract", "params": [42, 23], "id": 1}`,
			wantRequest: nil,
			wantError: &JsonRpcError{
				Code:    RpcInvalidRequest,
				Message: "invalid JSON-RPC version",
			},
		},
		{
			name:        "missing method",
			input:       `{"jsonrpc": "2.0", "params": [42, 23], "id": 1}`,
			wantRequest: nil,
			wantError: &JsonRpcError{
				Code:    RpcInvalidRequest,
				Message: "method is required",
			},
		},
		{
			name:        "invalid params type",
			input:       `{"jsonrpc": "2.0", "method": "subtract", "params": "invalid", "id": 1}`,
			wantRequest: nil,
			wantError: &JsonRpcError{
				Code:    RpcInvalidRequest,
				Message: "params must be an array or an object",
			},
		},
		{
			name:        "invalid key",
			input:       `{"jsonrpc": "2.0", "method": "subtract", "params": [42, 23], "id": 1, "invalid": true}`,
			wantRequest: nil,
			wantError: &JsonRpcError{
				Code:    RpcInvalidRequest,
				Message: "invalid key: invalid",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// we convert the input to a map[string]interface{}
			var rawJson map[string]interface{}
			if err := json.Unmarshal([]byte(tt.input), &rawJson); err != nil {
				if tt.wantError == nil {
					t.Errorf("failed to unmarshal input: %v", err)
					return
				}
				// we skip the test if we expect an error
				t.Skip()
			}
			gotRequest, _, gotError := ParseSimpleRequest(rawJson)

			if !reflect.DeepEqual(gotError, tt.wantError) {
				t.Errorf("ParseRequest() error = %#v, wantError %#v", gotError, tt.wantError)
				return
			}

			if !reflect.DeepEqual(gotRequest, tt.wantRequest) {
				t.Errorf("ParseRequest() = %#v, want %#v", gotRequest, tt.wantRequest)
			}
		})
	}
}

// TODO: we should convert to json object for comparison
// string comparison is not a good idea as the order of the keys is not guaranteed

func TestMarshalJsonRpcRequest(t *testing.T) {
	tests := []struct {
		name      string
		request   *JsonRpcRequest
		wantJson  string
		wantError bool
	}{
		{
			name: "valid request with positional params",
			request: &JsonRpcRequest{
				JsonRpcVersion: "2.0",
				Method:         "subtract",
				Params:         &JsonRpcParams{PositionalParams: []interface{}{float64(42), float64(23)}},
				Id:             &JsonRpcRequestId{Number: intPtr(1)},
			},
			wantJson:  `{"jsonrpc":"2.0","method":"subtract","params":[42,23],"id":1}`,
			wantError: false,
		},
		{
			name: "valid request with named params",
			request: &JsonRpcRequest{
				JsonRpcVersion: "2.0",
				Method:         "subtract",
				Params:         &JsonRpcParams{NamedParams: map[string]interface{}{"minuend": float64(42), "subtrahend": float64(23)}},
				Id:             &JsonRpcRequestId{String: stringPtr("1")},
			},
			wantJson:  `{"jsonrpc":"2.0","method":"subtract","params":{"minuend":42,"subtrahend":23},"id":"1"}`,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJson, err := MarshalJsonRpcRequest(tt.request)
			if (err != nil) != tt.wantError {
				t.Errorf("MarshalJsonRpcRequest() error = %v, wantError %v", err, tt.wantError)
				return
			}
			if string(gotJson) != tt.wantJson {
				t.Errorf("MarshalJsonRpcRequest() = %s, want %s", string(gotJson), tt.wantJson)
			}
		})
	}
}
