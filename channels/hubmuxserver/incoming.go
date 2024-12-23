package hubmuxserver

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (s *MuxSession) handleIncomingMessage(message transport.JsonRpcMessage) error {
	if message.Response != nil {
		response := message.Response
		if response.Error != nil {
			s.logger.Error("error in response", types.LogArg{
				"response":      fmt.Sprintf("%+v", response),
				"error_message": response.Error.Message,
				"error_code":    response.Error.Code,
				"error_data":    response.Error.Data,
			})
			return nil
		}
		switch message.Method {
		case mux.RpcRequestMethodCallTool:
			{
				toolsCallResult, err := mux.ParseJsonRpcResponseToolsCall(response)
				if err != nil {
					s.logger.Error("Failed to parse response", types.LogArg{
						"response": fmt.Sprintf("%+v", response),
						"error":    err,
					})
					return err
				}
				s.events.EventMuxResponseToolCall(toolsCallResult, response.Id)
			}
		default:
			s.logger.Error("received response message with unexpected method", types.LogArg{
				"method":   message.Method,
				"response": fmt.Sprintf("%+v", response),
				"id":       response.Id,
			})
		}
	} else if message.Request != nil {
		request := message.Request
		switch message.Method {
		case mux.RpcRequestMethodProxyRegister:
			{
				params, err := mux.ParseJsonRpcRequestProxyRegisterParams(request)
				if err != nil {
					s.logger.Error("Failed to parse request params", types.LogArg{
						"request": request,
						"method":  request.Method,
						"error":   err,
					})
					return err
				}
				// set the session information, required to
				// send the event to a specific proxy
				proxyId := params.ProxyId
				if proxyId == "" {
					s.logger.Error("missing proxy id", types.LogArg{
						"request": request,
						"method":  request.Method,
					})
					return fmt.Errorf("missing proxy id")
				}
				// we store the proxy id in the session
				s.proxyId = proxyId
				s.proxyName = params.ServerInfo.Name

				s.logger.Info("@@ Proxy register", types.LogArg{
					"proxyId":   s.proxyId,
					"proxyName": s.proxyName,
				})

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
					return err
				}
				s.logger.Info("Tools register", types.LogArg{
					"tools": params.Tools,
				})

				// send the event
				s.events.EventMuxRequestToolsRegister(s.proxyId, params, request.Id)
			}

		default:
			s.SendError(jsonrpc.RpcMethodNotFound, fmt.Sprintf("unknown method: %s", request.Method), request.Id)
		}
	} else {
		s.logger.Error("received message with unexpected nature", types.LogArg{
			"message": message,
		})
	}

	return nil

}
