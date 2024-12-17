package hubmuxserver

import (
	"context"

	"github.com/llmcontext/gomcp/eventbus"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type MuxSession struct {
	sessionId string
	transport *transport.JsonRpcTransport
	logger    types.Logger
	eventBus  *eventbus.EventBus
}

func NewMuxSession(sessionId string, tran types.Transport, logger types.Logger, eventBus *eventbus.EventBus) *MuxSession {
	jsonRpcTransport := transport.NewJsonRpcTransport(tran, "gomcp - proxy (mux)", logger)

	session := &MuxSession{
		sessionId: sessionId,
		transport: jsonRpcTransport,
		logger:    logger,
		eventBus:  eventBus,
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

func (s *MuxSession) SendRequest(request *jsonrpc.JsonRpcRequest) error {
	return s.transport.SendRequest(request)
}

func (s *MuxSession) SendResponse(response *jsonrpc.JsonRpcResponse) error {
	return s.transport.SendResponse(response)
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
		err := handleProxyRegister(s, request)
		if err != nil {
			s.logger.Error("Failed to handle proxy register", types.LogArg{
				"request": request,
				"method":  request.Method,
				"error":   err,
			})
			return
		}
	}
}

func (s *MuxSession) Close() {
	if s.transport != nil {
		s.transport.Close()
		s.transport = nil
	}
}

func (s *MuxSession) SessionId() string {
	return s.sessionId
}
