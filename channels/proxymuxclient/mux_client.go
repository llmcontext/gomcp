package proxymuxclient

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/llmcontext/gomcp/channels/proxy/events"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/transport/socket"
	"github.com/llmcontext/gomcp/types"
)

type ProxyMuxClient struct {
	transport  *transport.JsonRpcTransport
	logger     types.Logger
	events     events.Events
	muxAddress string
}

func NewProxyMuxClient(
	muxAddress string,
	events events.Events,
	logger types.Logger,
) *ProxyMuxClient {
	return &ProxyMuxClient{
		muxAddress: muxAddress,
		transport:  nil,
		logger:     logger,
		events:     events,
	}
}

func (c *ProxyMuxClient) Start(ctx context.Context) error {
	var err error
	errMuxChan := make(chan error, 1)

	isMuxServerReady := isMuxServerReady(c.muxAddress, c.logger)
	if !isMuxServerReady {
		return fmt.Errorf("mux server is not ready")
	}

	c.logger.Info("@@@ mux server is ready", types.LogArg{"step": 1})

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

	c.transport = muxJsonRpcTransport

	// for that specific transport, we send the event that the mux client is started
	// because the socket connection is already established
	c.events.EventMuxStarted()

	go func() {
		err = c.transport.Start(ctx, func(msg transport.JsonRpcMessage, jsonRpcTransport *transport.JsonRpcTransport) {
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
	c.transport.Close()
}

func (s *ProxyMuxClient) SendJsonRpcResponse(response interface{}, id *jsonrpc.JsonRpcRequestId) {
	s.transport.SendResponse(&jsonrpc.JsonRpcResponse{
		JsonRpcVersion: jsonrpc.JsonRpcVersion,
		Id:             id,
		Result:         response,
		Error:          nil,
	})
}

func (s *ProxyMuxClient) SendRequestWithMethodAndParams(method string, params interface{}) (*jsonrpc.JsonRpcRequestId, error) {
	return s.transport.SendRequestWithMethodAndParams(method, params)
}

func (s *ProxyMuxClient) SendError(code int, message string, id *jsonrpc.JsonRpcRequestId) {
	s.logger.Debug("JsonRpcError", types.LogArg{
		"code":    code,
		"message": message,
		"id":      id,
	})
	err := s.transport.SendError(code, message, id)
	if err != nil {
		s.logger.Error("failed to send error", types.LogArg{
			"error": err,
		})
	}
}

func isMuxServerReady(muxAddress string, logger types.Logger) bool {
	var maxAttempts = 60
	for i := 0; i < maxAttempts; i++ {
		logger.Info("waiting for mux server to be ready", types.LogArg{"attempts": i})
		conn, err := net.DialTimeout("tcp", muxAddress, 5*time.Second)
		if conn != nil {
			conn.Close()
		}
		if err == nil {
			return true
		}
		time.Sleep(2 * time.Second)
	}
	logger.Error("mux server is not ready", types.LogArg{"attempts": maxAttempts})
	return false
}
