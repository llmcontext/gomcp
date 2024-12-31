package mcpserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (m *McpServer) startProtocol(ctx context.Context, tran types.Transport) error {
	// create a new json rpc transport
	jsonRpcTransport := transport.NewJsonRpcTransport(tran, "mcp server", m.logger)
	m.jsonRpcTransport = jsonRpcTransport

	var err error

	errChan := make(chan error, 1)

	go func() {
		// Start the transport
		err := jsonRpcTransport.Start(ctx, func(message transport.JsonRpcMessage, jsonRpcTransport *transport.JsonRpcTransport) {
			err = m.handleIncomingMessage(ctx, message)
			if err != nil {
				m.logger.Error("failed to handle incoming message", types.LogArg{
					"error": err,
				})
			}
		})
		if err != nil {
			m.logger.Error("failed to start transport", types.LogArg{
				"error": err,
			})
		}
		errChan <- err
	}()

	select {
	case err := <-errChan:
		return err
	case <-ctx.Done():
		jsonRpcTransport.Close()
		return ctx.Err()
	}
}

func (m *McpServer) handleIncomingMessage(
	ctx context.Context,
	message transport.JsonRpcMessage,
) error {
	if message.Response != nil {
		response := message.Response
		if response.Error != nil {
			m.logger.Error("error in response", types.LogArg{
				"response":      fmt.Sprintf("%+v", response),
				"error_message": response.Error.Message,
				"error_code":    response.Error.Code,
				"error_data":    response.Error.Data,
			})
			return nil
		}
		switch message.Method {
		default:
			m.logger.Error("received message with unexpected method", types.LogArg{
				"method": message.Method,
				"c":      "p11h",
			})
		}
	} else if message.Request != nil {
		request := message.Request
		switch message.Method {
		case mcp.RpcRequestMethodInitialize:
			{
				parsed, err := mcp.ParseJsonRpcRequestInitialize(request)
				if err != nil {
					m.jsonRpcTransport.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
				}
				m.EventMcpRequestInitialize(parsed, request.Id)
			}
		case mcp.RpcNotificationMethodInitialized:
			m.EventMcpNotificationInitialized()
		case mcp.RpcRequestMethodToolsList:
			{
				parsed, err := mcp.ParseJsonRpcRequestToolsList(request)
				if err != nil {
					m.jsonRpcTransport.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
				}
				m.EventMcpRequestToolsList(parsed, request.Id)
			}
		case mcp.RpcRequestMethodToolsCall:
			{
				parsed, err := mcp.ParseJsonRpcRequestToolsCallParams(request.Params)
				if err != nil {
					m.jsonRpcTransport.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
					return nil
				}
				m.EventMcpRequestToolsCall(ctx, parsed, request.Id)
			}
		case mcp.RpcRequestMethodResourcesList:
			{
				parsed, err := mcp.ParseJsonRpcRequestResourcesList(request.Params)
				if err != nil {
					m.jsonRpcTransport.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
				}
				m.EventMcpRequestResourcesList(parsed, request.Id)
			}
		case mcp.RpcRequestMethodPromptsList:
			{
				parsed, err := mcp.ParseJsonRpcRequestPromptsList(request.Params)
				if err != nil {
					m.jsonRpcTransport.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
				}
				m.EventMcpRequestPromptsList(parsed, request.Id)
			}
		case mcp.RpcRequestMethodPromptsGet:
			{
				parsed, err := mcp.ParseJsonRpcRequestPromptsGet(request.Params)
				if err != nil {
					m.jsonRpcTransport.SendError(jsonrpc.RpcInvalidRequest, err.Error(), request.Id)
				}
				m.EventMcpRequestPromptsGet(parsed, request.Id)
			}
		case "ping":
			result := json.RawMessage(`{}`)
			m.jsonRpcTransport.SendJsonRpcResponse(result, request.Id)
		default:
			m.jsonRpcTransport.SendError(jsonrpc.RpcMethodNotFound, fmt.Sprintf("unknown method: %s", request.Method), request.Id)
		}
	} else {
		m.logger.Error("received message with unexpected nature", types.LogArg{
			"message": message,
		})
	}

	return nil
}

func (m *McpServer) EventMcpRequestInitialize(params *mcp.JsonRpcRequestInitializeParams, reqId *jsonrpc.JsonRpcRequestId) {
	// store client information
	if params.ProtocolVersion != mcp.ProtocolVersion {
		m.logger.Error("protocol version mismatch", types.LogArg{
			"expected": mcp.ProtocolVersion,
			"received": params.ProtocolVersion,
		})
		m.jsonRpcTransport.SendError(jsonrpc.RpcInvalidRequest, "protocol version mismatch", reqId)
		return
	}
	// we store the client information
	m.clientName = params.ClientInfo.Name
	m.clientVersion = params.ClientInfo.Version

	// prepare response
	response := mcp.JsonRpcResponseInitializeResult{
		ProtocolVersion: mcp.ProtocolVersion,
		Capabilities: mcp.ServerCapabilities{
			Tools: &mcp.ServerCapabilitiesTools{
				ListChanged: jsonrpc.BoolPtr(true),
			},
			Prompts: &mcp.ServerCapabilitiesPrompts{
				ListChanged: jsonrpc.BoolPtr(true),
			},
		},
		ServerInfo: mcp.ServerInfo{Name: m.serverName, Version: m.serverVersion},
	}
	m.jsonRpcTransport.SendJsonRpcResponse(&response, reqId)
}
func (m *McpServer) EventMcpNotificationInitialized() {
	// that's a notification, no response is needed
	m.isClientInitialized = true
}

func (m *McpServer) EventMcpRequestToolsList(params *mcp.JsonRpcRequestToolsListParams, reqId *jsonrpc.JsonRpcRequestId) {
	var response = mcp.JsonRpcResponseToolsListResult{
		Tools: make([]mcp.ToolDescription, 0, 10),
	}

	// TODO: create the list of tools

	m.jsonRpcTransport.SendJsonRpcResponse(&response, reqId)
}

func (m *McpServer) EventMcpRequestToolsCall(ctx context.Context, params *mcp.JsonRpcRequestToolsCallParams, reqId *jsonrpc.JsonRpcRequestId) {
	// TODO: we retrieve the tool and "run" it
}

func (m *McpServer) EventMcpRequestResourcesList(params *mcp.JsonRpcRequestResourcesListParams, reqId *jsonrpc.JsonRpcRequestId) {
	var response = mcp.JsonRpcResponseResourcesListResult{
		Resources: make([]mcp.ResourceDescription, 0),
	}

	m.jsonRpcTransport.SendJsonRpcResponse(&response, reqId)
}

func (m *McpServer) EventMcpRequestPromptsList(params *mcp.JsonRpcRequestPromptsListParams, reqId *jsonrpc.JsonRpcRequestId) {
	var response = mcp.JsonRpcResponsePromptsListResult{
		Prompts: make([]mcp.PromptDescription, 0),
	}

	// prompts := m.promptsRegistry.GetListOfPrompts()
	// for _, prompt := range prompts {
	// 	arguments := make([]mcp.PromptArgumentDescription, 0, len(prompt.Arguments))
	// 	for _, argument := range prompt.Arguments {
	// 		arguments = append(arguments, mcp.PromptArgumentDescription{
	// 			Name:        argument.Name,
	// 			Description: argument.Description,
	// 			Required:    argument.Required,
	// 		})
	// 	}
	// 	response.Prompts = append(response.Prompts, mcp.PromptDescription{
	// 		Name:        prompt.Name,
	// 		Description: prompt.Description,
	// 		Arguments:   arguments,
	// 	})
	// }

	m.jsonRpcTransport.SendJsonRpcResponse(&response, reqId)
}

func (m *McpServer) EventMcpRequestPromptsGet(params *mcp.JsonRpcRequestPromptsGetParams, reqId *jsonrpc.JsonRpcRequestId) {
	var templateArgs = map[string]string{}
	// copy the arguments, as strings
	for key, value := range params.Arguments {
		templateArgs[key] = fmt.Sprintf("%v", value)
	}
	// promptName := params.Name

	// response, err := s.promptsRegistry.GetPrompt(promptName, templateArgs)
	// if err != nil {
	// 	s.mcpServer.SendError(jsonrpc.RpcInvalidParams, fmt.Sprintf("prompt processing error: %s", err), reqId)
	// 	return
	// }

	response := json.RawMessage(`{ "result": "hello" }`)

	// marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		m.jsonRpcTransport.SendError(jsonrpc.RpcInternalError, "failed to marshal response", reqId)
	}
	jsonResponse := json.RawMessage(responseBytes)

	// we send the response
	m.jsonRpcTransport.SendJsonRpcResponse(&jsonResponse, reqId)
}
