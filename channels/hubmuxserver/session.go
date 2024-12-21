package hubmuxserver

import (
	"context"

	"github.com/google/uuid"
	"github.com/llmcontext/gomcp/channels/hub/events"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type MuxSession struct {
	sessionId string
	transport *transport.JsonRpcTransport
	logger    types.Logger
	proxyId   string
	proxyName string
	events    events.Events
}

func NewMuxSession(sessionId string, tran types.Transport, logger types.Logger, events events.Events) *MuxSession {
	jsonRpcTransport := transport.NewJsonRpcTransport(tran, "gomcp - proxy (mux)", logger)

	session := &MuxSession{
		sessionId: sessionId,
		transport: jsonRpcTransport,
		logger:    logger,
		events:    events,
	}

	return session
}

func (s *MuxSession) Start(ctx context.Context) error {
	errChan := make(chan error, 1)

	go func() {
		err := s.transport.Start(ctx, func(message transport.JsonRpcMessage, jsonRpcTransport *transport.JsonRpcTransport) {
			if message.Response != nil {
				s.onJsonRpcResponse(message.Response)
			} else if message.Request != nil {
				s.onJsonRpcRequest(message.Request)
			}
		})
		if err != nil {
			s.logger.Error("Failed to start JSON-RPC transport", types.LogArg{
				"error": err,
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

func (s *MuxSession) SetSessionInformation(proxyId string, serverName string) {
	s.proxyId = proxyId
	s.proxyName = serverName
}

func (s *MuxSession) SessionId() string {
	return s.sessionId
}

func (s *MuxSession) ProxyId() string {
	return s.proxyId
}

func (s *MuxSession) ProxyName() string {
	return s.proxyName
}

func (s *MuxSession) SendRequest(request *jsonrpc.JsonRpcRequest) error {
	return s.transport.SendRequest(request, "")
}

func (s *MuxSession) SendRequestWithExtraParam(request *jsonrpc.JsonRpcRequest, extraParam string) error {
	return s.transport.SendRequest(request, extraParam)
}

func (s *MuxSession) SendJsonRpcResponse(response interface{}, id *jsonrpc.JsonRpcRequestId) {
	s.transport.SendResponse(&jsonrpc.JsonRpcResponse{
		JsonRpcVersion: jsonrpc.JsonRpcVersion,
		Id:             id,
		Result:         response,
		Error:          nil,
	})
}

func (s *MuxSession) SendResponse(response *jsonrpc.JsonRpcResponse) error {
	return s.transport.SendResponse(response)
}

func (s *MuxSession) SendRequestWithMethodAndParams(method string, params interface{}, extraParam string) error {
	return s.transport.SendRequestWithMethodAndParams(method, params, extraParam)
}

func (s *MuxSession) onJsonRpcResponse(response *jsonrpc.JsonRpcResponse) {
	s.logger.Info("Response", types.LogArg{
		"response": response,
	})
}

func (s *MuxSession) onJsonRpcRequest(request *jsonrpc.JsonRpcRequest) {
	s.logger.Info("Request", types.LogArg{
		"request": request,
	})
	switch request.Method {
	case mux.RpcRequestMethodProxyRegister:
		{
			params, err := mux.ParseJsonRpcRequestProxyRegisterParams(request)
			if err != nil {
				s.logger.Error("Failed to parse request params", types.LogArg{
					"request": request,
					"method":  request.Method,
					"error":   err,
				})
				return
			}
			// set the session information, required to
			// send the event to a specific proxy
			// TODO: store in database
			proxyId := params.ProxyId
			if proxyId == "" {
				// we need to generate a new proxy id
				proxyId = uuid.New().String()
			}
			s.proxyId = proxyId
			s.proxyName = params.ServerInfo.Name

			// send the event
			s.events.EventMuxRequestProxyRegister(s.proxyId, params, request.Id)

		}
	case mux.RpcRequestMethodToolsRegister:
		{
			params, err := mux.ParseJsonRpcRequestToolsRegisterParams(request)
			if err != nil {
				s.logger.Error("Failed to parse request params", types.LogArg{
					"request": request,
					"method":  request.Method,
					"error":   err,
				})
				return
			}
			s.logger.Info("Tools register", types.LogArg{
				"tools": params.Tools,
			})

			// send the event
			s.events.EventMuxRequestToolsRegister(s.proxyId, params, request.Id)
		}
	default:
		s.logger.Error("Unknown method", types.LogArg{
			"method":  request.Method,
			"request": request,
		})
	}
}

func (s *MuxSession) Close() {
	if s.transport != nil {
		s.transport.Close()
		s.transport = nil
	}
}
