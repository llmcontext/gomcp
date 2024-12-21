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
	onStarted     func()
	// pendingRequests is a map of message id to pending request
	pendingRequests map[string]*pendingRequest
	name            string
}

type JsonRpcMessage struct {
	Method     string
	Request    *jsonrpc.JsonRpcRequest
	Response   *jsonrpc.JsonRpcResponse
	ExtraParam string
}

func (m *JsonRpcMessage) IsRequest() bool {
	return m.Request != nil
}

func (m *JsonRpcMessage) IsResponse() bool {
	return m.Response != nil
}

func (m *JsonRpcMessage) DebugInfo(transportName string) map[string]interface{} {
	info := make(map[string]interface{})
	info["transportName"] = transportName
	if m.IsRequest() {
		info["nature"] = "request"
		info["method"] = m.Request.Method
		info["id"] = jsonrpc.RequestIdToString(m.Request.Id)
		info["params"] = m.Request.Params.String()
	}
	if m.IsResponse() {
		info["nature"] = "response"
		if m.Response.Error != nil {
			info["error"] = fmt.Sprintf("%d: %s", m.Response.Error.Code, m.Response.Error.Message)
		}
		info["result"] = fmt.Sprintf("%v", m.Response.Result)
		info["id"] = jsonrpc.RequestIdToString(m.Response.Id)
	}
	return info
}

func NewJsonRpcTransport(transport types.Transport, name string, logger types.Logger) *JsonRpcTransport {
	return &JsonRpcTransport{
		transport:       transport,
		logger:          logger,
		lastRequestId:   0,
		pendingRequests: make(map[string]*pendingRequest),
		name:            name,
	}
}

func (m *JsonRpcTransport) Name() string {
	return m.name
}

func (m *JsonRpcTransport) OnStarted(callback func()) {
	m.onStarted = callback
}

func (t *JsonRpcTransport) Start(ctx context.Context, onMessage func(message JsonRpcMessage, jsonRpcTransport *JsonRpcTransport)) error {
	t.transport.OnMessage(func(message json.RawMessage) {
		// Check context before processing message
		select {
		case <-ctx.Done():
			return
		default:
			t.logger.Debug("message received", types.LogArg{
				"name":    t.name,
				"message": string(message),
			})
			// check the message nature
			nature, jsonRpcRawMessage, err := jsonrpc.CheckJsonMessage(message)
			if err != nil {
				t.logger.Error("error checking message", types.LogArg{
					"message": string(message),
					"error":   err,
					"name":    t.name,
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
						"name":  t.name,
					})
					return
				}
				onMessage(JsonRpcMessage{
					Request:  request,
					Method:   request.Method,
					Response: nil,
				}, t)

			case jsonrpc.MessageNatureResponse:
				response, _, rpcErr := jsonrpc.ParseJsonRpcResponse(jsonRpcRawMessage)

				if rpcErr != nil {
					t.logger.Error("error parsing response", types.LogArg{
						"error": rpcErr,
						"name":  t.name,
					})
					return
				}
				pendingRequestMethod, reqId := t.GetPendingRequest(response.Id)
				if pendingRequestMethod == "" {
					t.logger.Error("pending request method not found", types.LogArg{
						"requestId": jsonrpc.RequestIdToString(response.Id),
						"name":      t.name,
						"response":  response,
					})
					return
				}
				t.logger.Info("pending request method found", types.LogArg{
					"method": pendingRequestMethod,
					"name":   t.name,
					"id":     jsonrpc.RequestIdToString(reqId),
				})
				onMessage(JsonRpcMessage{
					Response: response,
					Method:   pendingRequestMethod,
				}, t)
			default:
				t.logger.Error("invalid message nature", types.LogArg{
					"nature": nature,
					"name":   t.name,
				})
				return
			}
		}
	})

	t.transport.OnClose(func() {
		t.logger.Info("transport closed", types.LogArg{
			"name": t.name,
		})
	})

	t.transport.OnError(func(err error) {
		t.logger.Error("transport error", types.LogArg{
			"error": err,
			"name":  t.name,
		})
	})

	t.transport.OnStarted(func() {
		if t.onStarted != nil {
			t.onStarted()
		}
	})

	errChan := make(chan error, 1)
	go func() {
		// Start the transport
		err := t.transport.Start(ctx)
		if err != nil {
			t.logger.Error("error starting transport", types.LogArg{
				"error": err,
				"name":  t.name,
			})
		}
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (t *JsonRpcTransport) SendRequestWithMethodAndParams(method string, params interface{}) (*jsonrpc.JsonRpcRequestId, error) {
	requestId := t.GetNextRequestId()
	request := buildJsonRpcRequestWithNamedParams(
		method, params, requestId)

	if request == nil {
		return nil, fmt.Errorf("failed to create initialize request")
	}

	return requestId, t.SendRequest(request)
}

func (t *JsonRpcTransport) SendResponseWithResults(reqId *jsonrpc.JsonRpcRequestId, result interface{}) error {
	response := &jsonrpc.JsonRpcResponse{
		JsonRpcVersion: jsonrpc.JsonRpcVersion,
		Id:             reqId,
		Result:         result,
		Error:          nil,
	}
	return t.SendResponse(response)
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

	t.logger.Info("sending request", types.LogArg{
		"method":  request.Method,
		"id":      jsonrpc.RequestIdToString(request.Id),
		"name":    t.name,
		"request": request,
	})

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

func (t *JsonRpcTransport) SendError(code int, message string, id *jsonrpc.JsonRpcRequestId) error {
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
		t.logger.Error("pending request not found", types.LogArg{
			"requestId": reqIdStr,
		})
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
