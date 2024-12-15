package transport

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/types"
)

type pendingRequest struct {
	method    string
	requestId *jsonrpc.JsonRpcRequestId
}

type JsonRpcTransport struct {
	transport     types.Transport
	logger        types.Logger
	lastRequestId int
	// pendingRequests is a map of message id to pending request
	pendingRequests map[string]*pendingRequest
}

type JsonRpcMessage struct {
	Request  *jsonrpc.JsonRpcRequest
	Response *jsonrpc.JsonRpcResponse
}

func NewJsonRpcTransport(transport types.Transport, logger types.Logger) *JsonRpcTransport {
	return &JsonRpcTransport{
		transport:       transport,
		logger:          logger,
		lastRequestId:   0,
		pendingRequests: make(map[string]*pendingRequest),
	}
}

func (t *JsonRpcTransport) Start(ctx context.Context, onMessage func(message JsonRpcMessage)) error {
	t.transport.OnMessage(func(message json.RawMessage) {
		t.logger.Debug("message received", types.LogArg{
			"message": message,
		})
		// check the message nature
		nature, jsonRpcRawMessage, err := jsonrpc.CheckJsonMessage(message)
		if err != nil {
			t.logger.Error("error checking message", types.LogArg{
				"error": err,
			})
			return
		}

		// the MCP protocol does not support batch requests
		switch nature {
		case jsonrpc.MessageNatureBatchRequest:
			t.logger.Error("batch requests not supported in protocol", types.LogArg{})
			return
		case jsonrpc.MessageNatureRequest:
			request, _, rpcErr := jsonrpc.ParseJsonRpcRequest(jsonRpcRawMessage)
			if rpcErr != nil {
				t.logger.Error("error parsing request", types.LogArg{
					"error": rpcErr,
				})
				return
			}
			onMessage(JsonRpcMessage{Request: request})

		case jsonrpc.MessageNatureResponse:
			response, _, rpcErr := jsonrpc.ParseJsonRpcResponse(jsonRpcRawMessage)
			if rpcErr != nil {
				t.logger.Error("error parsing response", types.LogArg{
					"error": rpcErr,
				})
				return
			}
			onMessage(JsonRpcMessage{Response: response})
		default:
			t.logger.Error("invalid message nature", types.LogArg{
				"nature": nature,
			})
			return
		}
	})

	t.transport.OnClose(func() {
		t.logger.Info("transport closed", types.LogArg{})
	})

	t.transport.OnError(func(err error) {
		t.logger.Error("transport error", types.LogArg{
			"error": err,
		})
	})

	return t.transport.Start(ctx)
}

func (t *JsonRpcTransport) SendRequestWithMethodAndParams(method string, params interface{}) error {

	request := buildJsonRpcRequestWithNamedParams(
		method, params, t.GetNextRequestId())

	if request == nil {
		return fmt.Errorf("failed to create initialize request")
	}

	return t.SendRequest(request)
}

func (t *JsonRpcTransport) SendRequest(request *jsonrpc.JsonRpcRequest) error {
	jsonMessage, err := jsonrpc.MarshalJsonRpcRequest(request)
	if err != nil {
		t.logger.Error("error marshalling message", types.LogArg{
			"error": err,
		})
		return err
	}

	// we store the request in the pending requests map
	// so we can match the response with the request
	if request.Id != nil {
		t.pendingRequests[jsonrpc.RequestIdToString(request.Id)] = &pendingRequest{
			method:    request.Method,
			requestId: request.Id,
		}
	}

	return t.transport.Send(jsonMessage)
}

func (t *JsonRpcTransport) SendResponse(response *jsonrpc.JsonRpcResponse) error {
	jsonMessage, err := jsonrpc.MarshalJsonRpcResponse(response)
	if err != nil {
		t.logger.Error("error marshalling message", types.LogArg{
			"error": err,
		})
		return err
	}
	return t.transport.Send(jsonMessage)
}

func (t *JsonRpcTransport) sendError(code int, message string, id *jsonrpc.JsonRpcRequestId) error {
	response := &jsonrpc.JsonRpcResponse{
		Error: &jsonrpc.JsonRpcError{
			Code:    code,
			Message: message,
		},
		Id: id,
	}

	jsonMessage, err := jsonrpc.MarshalJsonRpcResponse(response)
	if err != nil {
		t.logger.Error("error marshalling message", types.LogArg{
			"error": err,
		})
		return err
	}
	// send the message
	return t.transport.Send(jsonMessage)

}

func (t *JsonRpcTransport) Close() {
	t.transport.Close()
}

func (t *JsonRpcTransport) GetNextRequestId() *jsonrpc.JsonRpcRequestId {
	requestId := t.lastRequestId
	t.lastRequestId++
	return &jsonrpc.JsonRpcRequestId{
		Number: &requestId,
	}
}

// GetPendingRequest returns the method and message id of the pending request
// if the request is not found, it returns an empty string and nil
func (t *JsonRpcTransport) GetPendingRequest(reqId *jsonrpc.JsonRpcRequestId) (string, *jsonrpc.JsonRpcRequestId) {
	if reqId == nil {
		return "", nil
	}
	reqIdStr := jsonrpc.RequestIdToString(reqId)
	pendingRequest := t.pendingRequests[reqIdStr]
	if pendingRequest == nil {
		return "", nil
	}
	// we delete the pending request from the map
	delete(t.pendingRequests, reqIdStr)
	return pendingRequest.method, pendingRequest.requestId
}

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

func buildJsonRpcRequestWithNamedParams(method string, params interface{}, id *jsonrpc.JsonRpcRequestId) *jsonrpc.JsonRpcRequest {
	namedParams, err := structToMap(params)
	if err != nil {
		return nil
	}
	return &jsonrpc.JsonRpcRequest{
		JsonRpcVersion: jsonrpc.JsonRpcVersion,
		Method:         method,
		Params:         &jsonrpc.JsonRpcParams{NamedParams: namedParams},
		Id:             id,
	}
}
