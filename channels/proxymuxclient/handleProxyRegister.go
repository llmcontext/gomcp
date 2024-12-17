package proxymuxclient

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/types"
)

func (c *ProxyMuxClient) handleProxyRegisterResponse(response *jsonrpc.JsonRpcResponse) error {
	registerResponse, err := mux.ParseJsonRpcResponseProxyRegister(response)
	if err != nil {
		c.logger.Error("error in handleProxyRegisterResponse", types.LogArg{
			"error": err,
		})
		return err
	}

	c.logger.Info("proxy register response", types.LogArg{
		"sessionId":  registerResponse.SessionId,
		"proxyId":    registerResponse.ProxyId,
		"persistent": registerResponse.Persistent,
		"denied":     registerResponse.Denied,
	})

	return nil
}
