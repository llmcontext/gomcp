package hubmuxserver

import (
	"fmt"

	"github.com/google/uuid"
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
		default:
			s.logger.Error("received message with unexpected method", types.LogArg{
				"method": message.Method,
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
