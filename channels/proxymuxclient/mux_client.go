package proxymuxclient

import (
	"context"
	"fmt"

	"github.com/llmcontext/gomcp/channels/proxy/events"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/transport/socket"
	"github.com/llmcontext/gomcp/types"
)

type ProxyMuxClient struct {
	muxJsonRpcTransport *transport.JsonRpcTransport
	logger              types.Logger
	events              events.Events
	muxAddress          string
}

func NewProxyMuxClient(
	muxAddress string,
	events events.Events,
	logger types.Logger,
) *ProxyMuxClient {
	return &ProxyMuxClient{
		muxAddress:          muxAddress,
		muxJsonRpcTransport: nil,
		logger:              logger,
		events:              events,
	}
}

func (c *ProxyMuxClient) Start(ctx context.Context) error {
	var err error
	errMuxChan := make(chan error, 1)

	// start the mux client
	// create a transport for the mux client
	muxClientSocket := socket.NewSocketClient(c.muxAddress)

	// we try to start the mux client socket
	// let's get a transport for the mux client
	// if we fail, we return an error and stop the proxy client
	muxClientTransport, err := muxClientSocket.Start()
	if err != nil {
		c.logger.Error("error starting mux client socket", types.LogArg{"error": err})
		return err
	}

	if muxClientTransport == nil {
		c.logger.Error("error starting mux client socket", types.LogArg{"error": err})
		return fmt.Errorf("error creating mux transport")
	}

	// create the json rpc transport for the mux client
	muxJsonRpcTransport := transport.NewJsonRpcTransport(muxClientTransport, "proxy client - gomcp (mux)", c.logger)

	c.muxJsonRpcTransport = muxJsonRpcTransport
	go func() {
		err = c.muxJsonRpcTransport.Start(ctx, func(msg transport.JsonRpcMessage, jsonRpcTransport *transport.JsonRpcTransport) {
			c.logger.Debug("received message from mux", types.LogArg{
				"message":  msg,
				"method":   msg.Method,
				"request":  msg.Request,
				"response": msg.Response,
				"name":     jsonRpcTransport.Name(),
			})
			err = c.handleIncomingMessage(msg)
			if err != nil {
				c.logger.Error("error handling incoming message", types.LogArg{
					"error": err,
				})
				errMuxChan <- err
			}
		})
		if err != nil {
			c.logger.Error("failed to start mux transport", types.LogArg{
				"error": err,
			})
			errMuxChan <- err
		}
	}()

	select {
	case err := <-errMuxChan:
		c.Close()
		return err
	case <-ctx.Done():
		c.Close()
		return ctx.Err()
	}

}

func (c *ProxyMuxClient) Close() {
	c.muxJsonRpcTransport.Close()
}

func (c *ProxyMuxClient) SendToolCallResponse(toolsCallResult *mcp.JsonRpcResponseToolsCallResult, reqId *jsonrpc.JsonRpcRequestId, mcpReqId string) {
	params := mux.JsonRpcResponseToolsCallResult{
		Content:  toolsCallResult.Content,
		IsError:  toolsCallResult.IsError,
		McpReqId: mcpReqId,
	}
	c.muxJsonRpcTransport.SendResponseWithResults(reqId, params)
}

func (s *ProxyMuxClient) SendJsonRpcResponse(response interface{}, id *jsonrpc.JsonRpcRequestId) {
	s.muxJsonRpcTransport.SendResponse(&jsonrpc.JsonRpcResponse{
		JsonRpcVersion: jsonrpc.JsonRpcVersion,
		Id:             id,
		Result:         response,
		Error:          nil,
	})
}

func (s *ProxyMuxClient) SendRequestWithMethodAndParams(method string, params interface{}) {
	s.muxJsonRpcTransport.SendRequestWithMethodAndParams(method, params, "")
}

func (s *ProxyMuxClient) SendError(code int, message string, id *jsonrpc.JsonRpcRequestId) {
	s.logger.Debug("JsonRpcError", types.LogArg{
		"code":    code,
		"message": message,
		"id":      id,
	})
	err := s.muxJsonRpcTransport.SendError(code, message, id)
	if err != nil {
		s.logger.Error("failed to send error", types.LogArg{
			"error": err,
		})
	}
}
