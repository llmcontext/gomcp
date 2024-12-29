package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/registry"
	"github.com/llmcontext/gomcp/types"
	"github.com/llmcontext/gomcp/utils"
)

type ToolRpcHandler func(input json.RawMessage) (json.RawMessage, error)

type toolProviderPrepared struct {
	ToolProvider   *SdkToolProvider
	ToolDefinition *SdkToolDefinition
}

type ToolsRegistry struct {
	ToolProviders []*SdkToolProvider
	Tools         map[string]*toolProviderPrepared
	logger        types.Logger
}

func NewToolsRegistry(loadProxyTools bool, logger types.Logger) *ToolsRegistry {
	toolsRegistry := &ToolsRegistry{
		ToolProviders: []*SdkToolProvider{},
		Tools:         make(map[string]*toolProviderPrepared),
		logger:        logger,
	}
	// check if we need to load proxy tools
	if loadProxyTools {
		proxyTools := NewProxyTools()
		proxyTools.RegisterProxyTools(toolsRegistry)
	}
	return toolsRegistry
}

// TODO: need that?
func (r *ToolsRegistry) RegisterToolProvider(toolProvider *SdkToolProvider) error {
	r.ToolProviders = append(r.ToolProviders, toolProvider)
	r.logger.Info("registered tool provider", types.LogArg{
		"tool":            toolProvider.toolName,
		"contextTypeName": toolProvider.contextTypeName,
		"proxyId":         toolProvider.proxyId,
	})
	return nil
}

func (r *ToolsRegistry) RegisterProxyToolProvider(proxyId string, proxyName string) (*SdkToolProvider, error) {
	// check if the proxy tool provider is already registered
	for _, toolProvider := range r.ToolProviders {
		if toolProvider.proxyId == proxyId {
			return toolProvider, nil
		}
	}

	provider, err := newProxyToolProvider(proxyId, proxyName)
	if err != nil {
		return nil, err
	}
	r.ToolProviders = append(r.ToolProviders, provider)
	return provider, nil
}

func (r *ToolsRegistry) PrepareProxyToolProvider(toolProvider *SdkToolProvider) error {
	for _, toolDefinition := range toolProvider.toolDefinitions {
		r.Tools[toolDefinition.ToolName] = &toolProviderPrepared{
			ToolProvider:   toolProvider,
			ToolDefinition: toolDefinition,
		}
	}
	return nil
}

func (r *ToolsRegistry) Prepare(ctx context.Context) error {
	// let's prepare the different functions for each tool provider
	for _, toolProvider := range r.ToolProviders {
		if toolProvider.isDisabled {
			continue
		}
		// if the tool provider is a proxy, we don't need to prepare it
		// because it is already prepared by the proxy tools registry
		if toolProvider.proxyId != "" {
			continue
		}
		// for each tool definition, we prepare the function
		for _, toolDefinition := range toolProvider.toolDefinitions {
			// check that we don't already have a tool with this name
			if _, ok := r.Tools[toolDefinition.ToolName]; ok {
				return fmt.Errorf("tool %s already registered", toolDefinition.ToolName)
			}
			toolProviderPrepared := &toolProviderPrepared{
				ToolProvider:   toolProvider,
				ToolDefinition: toolDefinition,
			}
			r.Tools[toolDefinition.ToolName] = toolProviderPrepared
		}
	}

	// now, we can initialize the tool providers with their configuration
	// err := r.initializeProviders(ctx)
	// if err != nil {
	// 	return fmt.Errorf("error initializing tool providers: %w", err)
	// }

	return nil
}

func (r *ToolsRegistry) GetListOfTools() []*SdkToolDefinition {
	tools := make([]*SdkToolDefinition, 0, len(r.Tools))
	for _, tool := range r.Tools {
		tools = append(tools, tool.ToolDefinition)
	}
	return tools
}

func (r *ToolsRegistry) getTool(toolName string) (*SdkToolDefinition, *SdkToolProvider, error) {
	tool, ok := r.Tools[toolName]
	if !ok {
		return nil, nil, fmt.Errorf("tool %s not found", toolName)
	}
	return tool.ToolDefinition, tool.ToolProvider, nil
}

func (r *ToolsRegistry) IsProxyTool(toolName string) (bool, string, error) {
	_, toolProvider, err := r.getTool(toolName)
	if err != nil {
		return false, "", err
	}
	return toolProvider.proxyId != "", toolProvider.proxyId, nil
}

func (r *ToolsRegistry) CallTool(ctx context.Context, toolName string, toolArgs map[string]interface{}) (interface{}, error) {
	toolDefinition, toolProvider, err := r.getTool(toolName)
	if err != nil {
		return nil, err
	}

	// let's check if the arguments patch the schema
	err = utils.ValidateJsonSchemaWithObject(toolDefinition.InputSchema, toolArgs)
	if err != nil {
		return nil, err
	}

	// let's call the tool
	logger := types.NewSubLogger(r.logger, types.LogArg{
		"tool": toolProvider.toolName,
	})
	goCtx := types.ContextWithLogger(ctx, logger)

	// let's create the output
	output := registry.NewToolCallResult()

	_, callErr, err := utils.CallFunction(toolDefinition.ToolHandlerFunction, goCtx, toolProvider.toolContext, toolArgs, output)
	if err != nil {
		return nil, err
	}
	if callErr != nil {
		return nil, callErr
	}

	return output, nil
}
