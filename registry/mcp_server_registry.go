package registry

import (
	"context"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/types"
)

// abstract server interface

type PromptArgumentSchema struct {
	Name        string // the name of the argument
	Description string // the description of the argument
	Required    *bool  // true if the argument is required
}

type McpPromptLifecycle struct {
	Init    func(ctx context.Context, logger types.Logger) error
	Process func(ctx context.Context, params map[string]string, result types.PromptGetResult, errChan chan *jsonrpc.JsonRpcError)
	End     func(ctx context.Context, logger types.Logger) error
}

type McpToolLifecycle struct {
	IsInitialized bool
	Init          func(ctx context.Context, logger types.Logger) error
	Process       func(ctx context.Context, params map[string]interface{}, result types.ToolCallResult, logger types.Logger, errChan chan *jsonrpc.JsonRpcError) error
	End           func(ctx context.Context, logger types.Logger) error
}

type McpServerLifecycle struct {
	Init func(ctx context.Context, logger types.Logger) error
	End  func(ctx context.Context, logger types.Logger) error
}

type McpPromptDefinition struct {
	Name        string                 // the name of the prompt
	Description string                 // the description of the prompt
	Arguments   []PromptArgumentSchema // the arguments of the prompt
}

type McpToolDefinition struct {
	Name        string             // the name of the tool
	Description string             // the description of the tool
	InputSchema *jsonschema.Schema // A JSON Schema object defining the expected parameters for the tool, top object must be an object.
}

type McpPrompt struct {
	Definition *McpPromptDefinition
	Handler    *McpPromptLifecycle
}

type McpTool struct {
	Definition *McpToolDefinition
	Handler    *McpToolLifecycle
}

type McpServer struct {
	logger         types.Logger
	serverName     string
	serverVersion  string
	serverHandlers *McpServerLifecycle
	prompts        []McpPrompt
	tools          []McpTool
}

type McpServerRegistry struct {
	logger  types.Logger
	servers []*McpServer
}

func NewMcpServerRegistry(logger types.Logger) *McpServerRegistry {
	return &McpServerRegistry{
		logger:  logger,
		servers: make([]*McpServer, 0),
	}
}

func (r *McpServerRegistry) RegisterServer(serverName string, serverVersion string, handlers *McpServerLifecycle) (*McpServer, error) {
	server := &McpServer{
		serverName:     serverName,
		serverVersion:  serverVersion,
		serverHandlers: handlers,
		logger:         r.logger,
		prompts:        make([]McpPrompt, 0),
		tools:          make([]McpTool, 0),
	}
	r.logger.Debug("registry>server>RegisterServer", types.LogArg{
		"serverName":    serverName,
		"serverVersion": serverVersion,
	})
	r.servers = append(r.servers, server)
	return server, nil
}

func (s *McpServer) AddPrompt(prompt *McpPromptDefinition, handlers *McpPromptLifecycle) error {
	s.logger.Debug("registry>server>AddPrompt", types.LogArg{
		"promptName": prompt.Name,
	})
	for _, p := range s.prompts {
		if p.Definition.Name == prompt.Name {
			return fmt.Errorf("prompt already exists")
		}
	}
	s.prompts = append(s.prompts, McpPrompt{
		Definition: prompt,
		Handler:    handlers,
	})
	return nil
}

func (s *McpServer) AddTool(tool *McpToolDefinition, handlers *McpToolLifecycle) error {
	s.logger.Debug("registry>server>AddTool", types.LogArg{
		"toolName":        tool.Name,
		"toolDescription": tool.Description,
		"toolInputSchema": tool.InputSchema,
	})
	for _, t := range s.tools {
		if t.Definition.Name == tool.Name {
			return fmt.Errorf("tool already exists")
		}
	}
	s.tools = append(s.tools, McpTool{
		Definition: tool,
		Handler:    handlers,
	})
	return nil
}

func (r *McpServerRegistry) Init(ctx context.Context, logger types.Logger) error {
	for _, server := range r.servers {
		if server.serverHandlers.Init != nil {
			err := server.serverHandlers.Init(ctx, logger)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *McpServerRegistry) End(ctx context.Context, logger types.Logger) error {
	// end all the tools that have been initialized
	for _, server := range r.servers {
		for _, tool := range server.tools {
			if tool.Handler.IsInitialized {
				if tool.Handler.End != nil {
					err := tool.Handler.End(ctx, logger)
					if err != nil {
						logger.Error("error ending tool", types.LogArg{
							"toolName": tool.Definition.Name,
							"error":    err,
						})
					}
				}
			}
		}
	}

	// end all the servers
	for _, server := range r.servers {
		if server.serverHandlers.End != nil {
			err := server.serverHandlers.End(ctx, logger)
			if err != nil {
				logger.Error("error ending server", types.LogArg{
					"serverName": server.serverName,
					"error":      err,
				})
			}
		}
	}
	return nil
}

// return the list of tools from all the servers
func (r *McpServerRegistry) GetListOfTools() []McpToolDefinition {
	tools := make([]McpToolDefinition, 0)
	for _, server := range r.servers {
		for _, tool := range server.tools {
			tools = append(tools, *tool.Definition)
		}
	}
	return tools
}

func (r *McpServerRegistry) GetTool(toolName string) (*McpTool, error) {
	for _, server := range r.servers {
		for _, tool := range server.tools {
			if tool.Definition.Name == toolName {
				return &tool, nil
			}
		}
	}
	return nil, fmt.Errorf("tool not found")
}

func (t *McpTool) Run(
	ctx context.Context,
	arguments map[string]interface{},
	logger types.Logger,
) (types.ToolCallResult, *jsonrpc.JsonRpcError) {
	lifecycle := t.Handler
	if !lifecycle.IsInitialized {
		if lifecycle.Init != nil {
			err := lifecycle.Init(ctx, logger)
			if err != nil {
				return nil, &jsonrpc.JsonRpcError{
					Code:    jsonrpc.RpcInternalError,
					Message: err.Error(),
				}
			}
		}
		lifecycle.IsInitialized = true
	}

	// let's create the output
	output := NewToolCallResult()

	errChan := make(chan *jsonrpc.JsonRpcError, 1)
	go func() {
		lifecycle.Process(ctx, arguments, output, logger, errChan)
	}()

	// wait on context and errChan
	select {
	case err := <-errChan:
		if err != nil {
			return nil, err
		} else {
			return output, nil
		}
	case <-ctx.Done():
		return nil, &jsonrpc.JsonRpcError{
			Code:    jsonrpc.RpcInternalError,
			Message: ctx.Err().Error(),
		}
	}
}
