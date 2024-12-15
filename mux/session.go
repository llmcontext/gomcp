package mux

import (
	"context"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

type MuxSession struct {
	sessionId string
	transport *transport.JsonRpcTransport
	logger    types.Logger
}

func NewMuxSession(ctx context.Context, sessionId string, tran types.Transport, logger types.Logger) *MuxSession {
	jsonRpcTransport := transport.NewJsonRpcTransport(tran, logger)

	session := &MuxSession{
		sessionId: sessionId,
		transport: jsonRpcTransport,
		logger:    logger,
	}

	jsonRpcTransport.Start(ctx, func(message transport.JsonRpcMessage) {
		if message.Response != nil {
			session.onJsonRpcResponse(message.Response)
		} else if message.Request != nil {
			session.onJsonRpcRequest(message.Request)
		}
	})

	return session
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
}
