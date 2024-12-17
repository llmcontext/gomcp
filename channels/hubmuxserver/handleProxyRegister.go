package hubmuxserver

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/types"
)

func handleProxyRegister(s *MuxSession, request *jsonrpc.JsonRpcRequest) error {
	params, err := mux.ParseJsonRpcRequestProxyRegisterParams(request)
	if err != nil {
		s.logger.Error("Failed to parse request params", types.LogArg{
			"request": request,
			"method":  request.Method,
			"error":   err,
		})
		return err
	}
	s.logger.Info("Proxy registration", types.LogArg{
		"protocolVersion": params.ProtocolVersion,
		"proxyId":         params.ProxyId,
		"persistent":      params.Persistent,
		"proxy":           params.Proxy,
		"serverInfo":      params.ServerInfo,
	})
	// TODO: store in database

	// for now we accept all requests
	result := mux.JsonRpcResponseProxyRegisterResult{
		SessionId:  s.sessionId,
		ProxyId:    params.ProxyId,
		Persistent: params.Persistent,
		Denied:     false,
	}
	err = s.transport.SendResponseWithResults(request.Id, result)
	if err != nil {
		s.logger.Error("Failed to send response", types.LogArg{
			"error": err,
		})
		return err
	}

	return nil
}
