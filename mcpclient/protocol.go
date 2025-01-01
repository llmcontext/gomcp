package mcpclient

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (m *McpClient) handleIncomingMessage(message transport.JsonRpcMessage) error {
	if message.Response != nil {
		response := message.Response
		if response.Error != nil {
			m.logger.Error("error in response", types.LogArg{
				"method":        message.Method,
				"response":      fmt.Sprintf("%+v", response),
				"error_message": response.Error.Message,
				"error_code":    response.Error.Code,
				"error_data":    response.Error.Data,
			})
			m.jsonRpcTransport.SendError(response.Error.Code, response.Error.Message, response.Id)

			switch message.Method {
			case mcp.RpcRequestMethodToolsCall:
				// we forward the error
				m.notifications.OnToolCallResponse(nil, response.Id, response.Error)
			case mcp.RpcRequestMethodToolsList:
				// we forward the error
				m.notifications.OnToolsList(nil, response.Error)
			}
			return nil
		}
		switch message.Method {
		case mcp.RpcRequestMethodInitialize:
			{
				initializeResponse, err := mcp.ParseJsonRpcResponseInitialize(response)
				if err != nil {
					m.logger.Error("error in handleMcpInitializeResponse", types.LogArg{
						"error": err,
					})
					return nil
				}
				m.logger.Info("init response", types.LogArg{
					"name":    initializeResponse.ServerInfo.Name,
					"version": initializeResponse.ServerInfo.Version,
				})
				m.EventMcpResponseInitialize(initializeResponse)
			}
		case mcp.RpcRequestMethodToolsList:
			{
				// the MCP server sent its tools list
				toolsListResponse, err := mcp.ParseJsonRpcResponseToolsList(response)
				if err != nil {
					m.logger.Error("error in handleMcpToolsListResponse", types.LogArg{
						"error": err,
					})
					return nil
				}
				// we forward the tools list
				m.notifications.OnToolsList(toolsListResponse, nil)
			}
		case mcp.RpcRequestMethodToolsCall:
			{
				toolsCallResult, err := mcp.ParseJsonRpcResponseToolsCall(response)
				if err != nil {
					m.logger.Error("error parsing tools call params", types.LogArg{
						"error": err,
					})
					return nil
				}

				m.logger.Info("tools call result", types.LogArg{
					"content": toolsCallResult.Content,
					"isError": toolsCallResult.IsError,
				})

				// we forward the response
				m.notifications.OnToolCallResponse(toolsCallResult, response.Id, nil)
			}

		default:
			m.logger.Error("received message with unexpected method", types.LogArg{
				"method": message.Method,
				"c":      "4cdu",
			})
		}
	} else if message.Request != nil {
		request := message.Request
		switch message.Method {
		case mcp.RpcNotificationMethodResourcesUpdated:
			{
				resourcesUpdated, err := mcp.ParseJsonRpcNotificationResourcesUpdatedParams(request.Params)
				if err != nil {
					m.logger.Error("error parsing resources updated", types.LogArg{
						"error": err,
					})
					return nil
				}
				m.EventMcpNotificationResourcesUpdated(resourcesUpdated)
			}
		case mcp.RpcNotificationMethodResourcesListChanged:
			{
				m.EventMcpNotificationResourcesListChanged()
			}
		default:
			m.logger.Error("received message with unexpected method", types.LogArg{
				"method":  message.Method,
				"request": request,
				"c":       "cjp1",
			})
		}
	} else {
		m.logger.Error("received message with unexpected nature", types.LogArg{
			"message": message,
		})
	}

	return nil
}

func (m *McpClient) EventMcpStarted() {
	// as soon as the MCP server is started, we send an initialize request
	// to the MCP server
	params := mcp.JsonRpcRequestInitializeParams{
		ProtocolVersion: mcp.ProtocolVersion,
		Capabilities:    mcp.ClientCapabilities{},
		ClientInfo: mcp.ClientInfo{
			Name:    m.clientName,
			Version: m.clientVersion,
		},
	}
	m.jsonRpcTransport.SendRequestWithMethodAndParams(mcp.RpcRequestMethodInitialize, params)

}

func (m *McpClient) EventMcpResponseInitialize(resp *mcp.JsonRpcResponseInitializeResult) {
	m.logger.Debug("MCP Server initialized", types.LogArg{
		"name":    resp.ServerInfo.Name,
		"version": resp.ServerInfo.Version,
	})

	// we forward the server information
	m.notifications.OnServerInformation(resp.ServerInfo.Name, resp.ServerInfo.Version)

	// we send the "notifications/initialized" notification
	m.jsonRpcTransport.SendNotification(mcp.RpcNotificationMethodInitialized)

	// we send the "tools/list" request
	m.jsonRpcTransport.SendRequestWithMethodAndParams(
		mcp.RpcRequestMethodToolsList, mcp.JsonRpcRequestToolsListParams{})
}

// this is a tool call from the hub
func (m *McpClient) EventMuxRequestToolCall(toolName string, args map[string]interface{}) (*jsonrpc.JsonRpcRequestId, error) {
	m.logger.Info("EventMuxRequestToolCall", types.LogArg{
		"name": toolName,
		"args": args,
	})

	req := mcp.JsonRpcRequestToolsCallParams{
		Name:      toolName,
		Arguments: args,
	}

	mcpReqId, err := m.jsonRpcTransport.SendRequestWithMethodAndParams(mcp.RpcRequestMethodToolsCall, req)
	if err != nil {
		m.logger.Error("failed to send request to mcp client", types.LogArg{"error": err})
		return nil, err
	}

	return mcpReqId, nil
}

func (m *McpClient) EventMcpNotificationResourcesListChanged() {
	m.logger.Info("event mcp notification resources list changed", types.LogArg{})
}

func (m *McpClient) EventMcpNotificationResourcesUpdated(resourcesUpdated *mcp.JsonRpcNotificationResourcesUpdatedParams) {
	m.logger.Info("event mcp notification resources updated", types.LogArg{
		"uri": resourcesUpdated.Uri,
	})
}
