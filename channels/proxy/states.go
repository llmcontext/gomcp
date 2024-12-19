package proxy

import (
	"github.com/llmcontext/gomcp/channels"
	"github.com/llmcontext/gomcp/channels/proxy/events"
	"github.com/llmcontext/gomcp/channels/proxymcpclient"
	"github.com/llmcontext/gomcp/channels/proxymuxclient"
	"github.com/llmcontext/gomcp/protocol/mcp"
	"github.com/llmcontext/gomcp/types"
)

type StateManager struct {
	logger    types.Logger
	options   *channels.ProxiedMcpServerDescription
	muxClient *proxymuxclient.ProxyMuxClient
	mcpClient *proxymcpclient.ProxyMcpClient

	// serverInfo is the info about the MCP server we are connected to
	serverInfo mcp.ServerInfo
}

func NewStateManager(options *channels.ProxiedMcpServerDescription, logger types.Logger) *StateManager {
	return &StateManager{
		options:    options,
		logger:     logger,
		serverInfo: mcp.ServerInfo{},
	}
}

func (s *StateManager) Stop(err error) {
	s.logger.Info("stopping state manager", types.LogArg{"error": err})
}

func (s *StateManager) SetMuxClient(muxClient *proxymuxclient.ProxyMuxClient) {
	s.muxClient = muxClient
}

func (s *StateManager) SetProxyClient(proxyClient *proxymcpclient.ProxyMcpClient) {
	s.mcpClient = proxyClient
}

func (s *StateManager) AsEventsProcessor() events.EventsProcessor {
	return s
}

func (s *StateManager) EventMcpStared() {
	s.mcpClient.SendInitializeRequest()
}

func (s *StateManager) EventMcpInitializeResponse(resp *mcp.JsonRpcResponseInitializeResult) {
	s.logger.Info("event mcp initialize response", types.LogArg{
		"name":    resp.ServerInfo.Name,
		"version": resp.ServerInfo.Version,
	})

	// we update the server information
	s.serverInfo.Name = resp.ServerInfo.Name
	s.serverInfo.Version = resp.ServerInfo.Version

	// we send the "notifications/initialized" notification
	s.mcpClient.SendInitializedNotification()

	// we send the "tools/list" request
	s.mcpClient.SendToolsListRequest()
}

func (s *StateManager) EventMcpToolsListResponse(resp *mcp.JsonRpcResponseToolsListResult) {
	s.logger.Info("event mcp tools list response", types.LogArg{
		"tools": resp.Tools,
	})

	// we can now report the tools list to the mux server
	s.muxClient.SendProxyRegistrationRequest(s.options, s.serverInfo, resp.Tools)

}
