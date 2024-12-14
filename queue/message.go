package queue

import "github.com/llmcontext/gomcp/jsonrpc"

// this is the message used to communicate through channels
// between the different go routines
type ProtocolMessageType string

const (
	ProtocolMessageTypeRequest  ProtocolMessageType = "mcp"
	ProtocolMessageTypeResponse ProtocolMessageType = "mux"
)

type ProtocolMessage struct {
	Protocol ProtocolMessageType      `json:"protocol"`
	Request  *jsonrpc.JsonRpcRequest  `json:"request"`
	Response *jsonrpc.JsonRpcResponse `json:"response"`
}

func NewProtocolMessageRequest(protocol ProtocolMessageType, request *jsonrpc.JsonRpcRequest) *ProtocolMessage {
	return &ProtocolMessage{
		Protocol: protocol,
		Request:  request,
		Response: nil,
	}
}

func NewProtocolMessageResponse(protocol ProtocolMessageType, response *jsonrpc.JsonRpcResponse) *ProtocolMessage {
	return &ProtocolMessage{
		Protocol: protocol,
		Request:  nil,
		Response: response,
	}
}
