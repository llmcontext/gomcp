package hubmuxserver

import (
	"github.com/google/uuid"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/types"
)

func (s *MuxSession) handleProxyRegister(request *jsonrpc.JsonRpcRequest) error {
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
	proxyId := params.ProxyId
	if proxyId == "" {
		// we need to generate a new proxy id
		proxyId = uuid.New().String()
	}
	// we need to store the proxy id in the session
	s.proxyId = proxyId
	s.proxyName = params.ServerInfo.Name

	// for now we accept all requests
	result := mux.JsonRpcResponseProxyRegisterResult{
		SessionId:  s.sessionId,
		ProxyId:    proxyId,
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
