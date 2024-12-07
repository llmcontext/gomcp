package jsonrpc

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var debugTestResponse = false

func TestMarshalJsonRpcResponse(t *testing.T) {
	t.Run("success with string result", func(t *testing.T) {
		result := json.RawMessage(`"hello"`)
		strId := "123"
		response := &JsonRpcResponse{
			Id:     &JsonRpcRequestId{String: &strId},
			Result: &result,
		}

		data, err := MarshalJsonRpcResponse(response)
		assert.NoError(t, err)

		var parsed RawJsonRpcResponseMessage
		err = json.Unmarshal(data, &parsed)
		if debugTestResponse {
			fmt.Printf("[%s]: %s\n", t.Name(), string(data))
		}
		assert.NoError(t, err)
		assert.Equal(t, JsonRpcVersion, parsed.JsonRpcVersion)
		assert.Equal(t, json.RawMessage(`"123"`), parsed.Id)
		assert.Equal(t, &result, parsed.Result)
	})

	t.Run("success with numeric id", func(t *testing.T) {
		result := json.RawMessage(`42`)
		numId := 1
		response := &JsonRpcResponse{
			Id:     &JsonRpcRequestId{Number: &numId},
			Result: &result,
		}

		data, err := MarshalJsonRpcResponse(response)
		if debugTestResponse {
			fmt.Printf("[%s]: %s\n", t.Name(), string(data))
		}
		assert.NoError(t, err)

		var parsed RawJsonRpcResponseMessage
		err = json.Unmarshal(data, &parsed)
		assert.NoError(t, err)
		assert.Equal(t, json.RawMessage(`1`), parsed.Id)
	})

	t.Run("error response", func(t *testing.T) {
		strId := "123"
		response := &JsonRpcResponse{
			Id: &JsonRpcRequestId{String: &strId},
			Error: &JsonRpcError{
				Code:    -32600,
				Message: "Invalid Request",
			},
		}

		data, err := MarshalJsonRpcResponse(response)
		if debugTestResponse {
			fmt.Printf("[%s]: %s\n", t.Name(), string(data))
		}
		assert.NoError(t, err)

		var parsed RawJsonRpcErrorMessage
		err = json.Unmarshal(data, &parsed)
		assert.NoError(t, err)
		assert.Equal(t, JsonRpcVersion, parsed.JsonRpcVersion)
		assert.Equal(t, json.RawMessage(`"123"`), parsed.Id)
		assert.Equal(t, -32600, parsed.Error.Code)
		assert.Equal(t, "Invalid Request", parsed.Error.Message)
	})

	t.Run("null id", func(t *testing.T) {
		result := json.RawMessage(`"hello"`)
		response := &JsonRpcResponse{
			Id:     nil,
			Result: &result,
		}

		data, err := MarshalJsonRpcResponse(response)
		if debugTestResponse {
			fmt.Printf("[%s]: %s\n", t.Name(), string(data))
		}
		assert.NoError(t, err)

		var parsed RawJsonRpcResponseMessage
		err = json.Unmarshal(data, &parsed)
		if debugTestResponse {
			fmt.Printf("[%s]: %#v\n", t.Name(), parsed)
		}
		assert.NoError(t, err)
		// check that parsedId is the byte array for null
		assert.Equal(t, json.RawMessage(`null`), parsed.Id)
	})

	t.Run("error: both result and error set", func(t *testing.T) {
		result := json.RawMessage(`"hello"`)
		response := &JsonRpcResponse{
			Result: &result,
			Error: &JsonRpcError{
				Code:    -32600,
				Message: "Invalid Request",
			},
		}

		_, err := MarshalJsonRpcResponse(response)
		assert.Error(t, err)
		assert.Equal(t, "error and result cannot be set at the same time", err.Error())
	})

	t.Run("error: neither result nor error set", func(t *testing.T) {
		response := &JsonRpcResponse{
			Id: nil,
		}

		_, err := MarshalJsonRpcResponse(response)
		assert.Error(t, err)
		assert.Equal(t, "error or result is required", err.Error())
	})
}
