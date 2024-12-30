package proxies

import (
	"encoding/json"
	"fmt"

	"github.com/invopop/jsonschema"
	"github.com/llmcontext/gomcp/registry"
)

type ProxyTool struct {
	proxyId    string
	definition *ProxyToolDefinition
}

func NewProxyTool(
	proxyId string,
	definition *ProxyToolDefinition,
) *ProxyTool {
	return &ProxyTool{
		proxyId:    proxyId,
		definition: definition,
	}
}

func (t *ProxyTool) register(mcpServer *registry.McpServer) error {
	var schema *jsonschema.Schema
	switch s := t.definition.InputSchema.(type) {
	case *jsonschema.Schema:
		schema = s
	case map[string]interface{}:
		schema = &jsonschema.Schema{}
		// Unmarshal the map into the schema
		if err := mapToStruct(s, schema); err != nil {
			return fmt.Errorf("invalid schema format: %v", err)
		}
	default:
		return fmt.Errorf("inputSchema must be either *jsonschema.Schema or map[string]interface{}")
	}

	var tool = registry.McpToolDefinition{
		Name:        t.definition.Name,
		Description: t.definition.Description,
		InputSchema: schema,
	}

	// TODO: implement those lifecycle methods
	var handlers = registry.McpToolLifecycle{
		Init:    nil,
		Process: nil,
		End:     nil,
	}
	err := mcpServer.AddTool(&tool, &handlers)
	if err != nil {
		return err
	}
	return nil
}

func mapToStruct(input map[string]interface{}, output interface{}) error {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, output)
}
