package mcpClient

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/protocol/mux"
)

func (c *MCPProxyClient) handleToolsListResponse(response *jsonrpc.JsonRpcResponse) {
	toolsListResponse, err := mcp.ParseJsonRpcResponseToolsList(response)
	if err != nil {
		c.logger.Error(fmt.Sprintf("error in handleToolsListResponse: %+v\n", err))
		return
	}

	// display all the tool names and descriptions
	for _, tool := range toolsListResponse.Tools {
		c.logger.Info(fmt.Sprintf("tool: %s - %s\n", tool.Name, tool.Description))
		c.logger.Info(fmt.Sprintf("inputSchema: %+#v\n", tool.InputSchema))
		// we store the tools description
		c.tools = append(c.tools, tool)
	}

	// we can now report the tools list to the mux server
	c.sendProxyRegistrationRequest()
}

func (c *MCPProxyClient) sendProxyRegistrationRequest() {
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

	c.logger.Info(fmt.Sprintf("sending proxy registration request: %+#v\n", params))
	c.muxClient.SendRequest(mux.RpcRequestMethodProxyRegister, params)
}
