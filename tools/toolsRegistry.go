package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/llmcontext/gomcp/types"
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
	return toolsRegistry
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
