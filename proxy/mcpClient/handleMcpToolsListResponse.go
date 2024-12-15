package mcpClient

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/transport"
	"github.com/llmcontext/gomcp/types"
)

func (c *MCPProxyClient) handleMcpToolsListResponse(
	response *jsonrpc.JsonRpcResponse,
	transport *transport.JsonRpcTransport) {
	toolsListResponse, err := mcp.ParseJsonRpcResponseToolsList(response)
	if err != nil {
		c.logger.Error("error in handleMcpToolsListResponse", types.LogArg{
			"error": err,
		})
		return
	}

	// display all the tool names and descriptions
	for _, tool := range toolsListResponse.Tools {
		c.logger.Info("tool", types.LogArg{
			"name":        tool.Name,
			"description": tool.Description,
		})
		// we store the tools description
		c.tools = append(c.tools, tool)
	}

	// we can now report the tools list to the mux server
	c.sendProxyRegistrationRequest(transport)
}

func (c *MCPProxyClient) sendProxyRegistrationRequest(transport *transport.JsonRpcTransport) {
	params := mux.JsonRpcRequestProxyRegisterParams{
		ProtocolVersion: mux.MuxProtocolVersion,
		Proxy: mux.ProxyDescription{
			WorkingDirectory: c.options.CurrentWorkingDirectory,
			Command:          c.options.ProgramName,
			Args:             c.options.ProgramArgs,
		},
		ServerInfo: mux.ServerInfo{
			Name:    c.serverInfo.Name,
			Version: c.serverInfo.Version,
		},
		Tools: []mux.ToolDescription{},
	}
	// we add the tools to the request
	for _, tool := range c.tools {
		params.Tools = append(params.Tools, mux.ToolDescription{
			Name:        tool.Name,
			Description: tool.Description,
			InputSchema: tool.InputSchema,
		})
	}

	c.logger.Info("sending proxy registration request", types.LogArg{
		"params": params,
	})
	err := transport.SendRequestWithMethodAndParams(mux.RpcRequestMethodMuxInitialize, params)
	if err != nil {
		c.logger.Error("error sending proxy registration request", types.LogArg{
			"error": err,
		})
	}
}
