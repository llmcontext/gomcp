package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/types"
)

type Queue struct {
	// transport carrying the raw messages
	transport types.Transport
	protocol  ProtocolMessageType

	// callback for when an error occurs
	errorCallback func(error)

	// the protocol message channels for the processing go routines
	protocolMessageFromProcessingChan chan *ProtocolMessage
	protocolMessageToProcessingChan   chan *ProtocolMessage
}

// a queue provides the communication channels between the transports and
// the processing go routines
// the transport with send/receive jsonRawMessage and the processing go routines
// send/receive protocol.ProtocolMessage
func NewQueue(protocol ProtocolMessageType, transport types.Transport) *Queue {
	return &Queue{
		transport:                         transport,
		protocol:                          protocol,
		protocolMessageFromProcessingChan: make(chan *ProtocolMessage, 100),
		protocolMessageToProcessingChan:   make(chan *ProtocolMessage, 100),
		errorCallback:                     nil,
	}
}

func (q *Queue) OnError(callback func(error)) {
	q.errorCallback = callback
}

func (q *Queue) Start(ctx context.Context) {
	transport := q.transport

	transport.OnMessage(func(msg json.RawMessage) {
		// check the message nature
		nature, jsonRpcRawMessage, err := jsonrpc.CheckJsonMessage(msg)
		if err != nil {
			if q.errorCallback != nil {
				q.errorCallback(err)
			} else {
				fmt.Printf("[queue] error: %s\n", err)
			}
			return
		}

		// the MCP protocol does not support batch requests
		switch nature {
		case jsonrpc.MessageNatureBatchRequest:
			if q.errorCallback != nil {
				q.errorCallback(fmt.Errorf("batch requests not supported in protocol"))
			} else {
				fmt.Printf("[queue] batch requests not supported in protocol\n")
			}
			return
		case jsonrpc.MessageNatureRequest:
			request, _, rpcErr := jsonrpc.ParseJsonRpcRequest(jsonRpcRawMessage)
			if rpcErr != nil {
				if q.errorCallback != nil {
					q.errorCallback(fmt.Errorf("invalid request: %v", rpcErr))
				} else {
					fmt.Printf("[queue] error: %v\n", rpcErr)
				}
				return
			}
			q.protocolMessageToProcessingChan <- NewProtocolMessageRequest(q.protocol, request)

		case jsonrpc.MessageNatureResponse:
			response, _, rpcErr := jsonrpc.ParseJsonRpcResponse(jsonRpcRawMessage)
			if rpcErr != nil {
				if q.errorCallback != nil {
					q.errorCallback(fmt.Errorf("invalid response: %v", rpcErr))
				} else {
					fmt.Printf("[queue] error: %v\n", rpcErr)
				}
				return
			}
			q.protocolMessageToProcessingChan <- NewProtocolMessageResponse(q.protocol, response)
		default:
			if q.errorCallback != nil {
				q.errorCallback(fmt.Errorf("invalid message nature: %d", nature))
			} else {
				fmt.Printf("[queue] invalid message nature: %d\n", nature)
			}
			return
		}

	})

	// Set up error handler
	transport.OnError(func(err error) {
		if q.errorCallback != nil {
			q.errorCallback(err)
		} else {
			fmt.Printf("[queue] transport error: %s\n", err)
		}
	})

	// Start the transport
	if err := transport.Start(ctx); err != nil {
		if q.errorCallback != nil {
			q.errorCallback(err)
		} else {
			fmt.Printf("[queue] failed to start transport: %s\n", err)
		}
	}

	// wait on ctx to be done but also on the channel protocolMessageFromProcessingChan

	select {
	case <-ctx.Done():
		transport.Close()
		// closing the channels
		close(q.protocolMessageToProcessingChan)
		close(q.protocolMessageFromProcessingChan)
	case msg := <-q.protocolMessageFromProcessingChan:
		// we received a message from the processing go routine
		// we need to send it to the transport
		if msg.Request != nil {
			json, err := jsonrpc.MarshalJsonRpcRequest(msg.Request)
			if err != nil {
				if q.errorCallback != nil {
					q.errorCallback(fmt.Errorf("failed to marshal request: %v", err))
				} else {
					fmt.Printf("[queue] failed to marshal request: %v\n", err)
				}
			}
			err = transport.Send(json)
			if err != nil {
				if q.errorCallback != nil {
					q.errorCallback(fmt.Errorf("failed to send request: %v", err))
				} else {
					fmt.Printf("[queue] failed to send request: %v\n", err)
				}
			}
		} else if msg.Response != nil {
			json, err := jsonrpc.MarshalJsonRpcResponse(msg.Response)
			if err != nil {
				if q.errorCallback != nil {
					q.errorCallback(fmt.Errorf("failed to marshal response: %v", err))
				} else {
					fmt.Printf("[queue] failed to marshal response: %v\n", err)
				}
			}
			err = transport.Send(json)
			if err != nil {
				if q.errorCallback != nil {
					q.errorCallback(fmt.Errorf("failed to send response: %v", err))
				} else {
					fmt.Printf("[queue] failed to send response: %v\n", err)
				}
			}
		}
	}
}
