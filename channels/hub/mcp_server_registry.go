package hub

import (
	"fmt"

	"github.com/llmcontext/gomcp/mcp"
)

type McpPrompt struct {
	Prompt  *mcp.McpPrompt
	Handler mcp.McpPromptLifecycle
}

type McpTool struct {
	Tool    *mcp.McpTool
	Handler mcp.McpToolLifecycle
}

type McpServer struct {
	serverName     string
	serverVersion  string
	serverHandlers mcp.McpServerLifecycle
	prompts        []McpPrompt
	tools          []McpTool
}

type McpServerRegistry struct {
	servers []*McpServer
}

func NewMcpServerRegistry() mcp.McpServerRegistry {
	return &McpServerRegistry{
		servers: make([]*McpServer, 0),
	}
}

func (r *McpServerRegistry) RegisterServer(serverName string, serverVersion string, handlers mcp.McpServerLifecycle) (mcp.McpServer, error) {
	server := &McpServer{
		serverName:     serverName,
		serverVersion:  serverVersion,
		serverHandlers: handlers,
		prompts:        make([]McpPrompt, 0),
		tools:          make([]McpTool, 0),
	}
	r.servers = append(r.servers, server)
	return server, nil
}

func (s *McpServer) AddPrompt(prompt *mcp.McpPrompt, handlers mcp.McpPromptLifecycle) error {
	for _, p := range s.prompts {
		if p.Prompt.Name == prompt.Name {
			return fmt.Errorf("prompt already exists")
		}
	}
	s.prompts = append(s.prompts, McpPrompt{
		Prompt:  prompt,
		Handler: handlers,
	})
	return nil
}

func (s *McpServer) AddTool(tool *mcp.McpTool, handlers mcp.McpToolLifecycle) error {
	for _, t := range s.tools {
		if t.Tool.Name == tool.Name {
			return fmt.Errorf("tool already exists")
		}
	}
	s.tools = append(s.tools, McpTool{
		Tool:    tool,
		Handler: handlers,
	})
	return nil
}
