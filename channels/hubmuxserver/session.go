package hubmuxserver

import (
	"context"

	"github.com/llmcontext/gomcp/channels/hub/events"
	"github.com/llmcontext/gomcp/jsonrpc"
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
			err := s.handleIncomingMessage(message)
			if err != nil {
				s.logger.Error("Failed to handle incoming message", types.LogArg{
					"error": err,
				})
			}
		})
		if err != nil {
			s.logger.Error("End of MuxSession", types.LogArg{
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

func (s *MuxSession) SendJsonRpcResponse(response interface{}, id *jsonrpc.JsonRpcRequestId) {
	s.transport.SendResponse(&jsonrpc.JsonRpcResponse{
		JsonRpcVersion: jsonrpc.JsonRpcVersion,
		Id:             id,
		Result:         response,
		Error:          nil,
	})
}

func (s *MuxSession) SendRequestWithMethodAndParams(method string, params interface{}) (*jsonrpc.JsonRpcRequestId, error) {
	return s.transport.SendRequestWithMethodAndParams(method, params)
}

func (s *MuxSession) SendError(code int, message string, id *jsonrpc.JsonRpcRequestId) {
	s.transport.SendError(code, message, id)
}

func (s *MuxSession) Close() {
	if s.transport != nil {
		s.transport.Close()
		s.transport = nil
	}
}
