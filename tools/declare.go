package tools

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/invopop/jsonschema"
)

type SdkToolDefinition struct {
	ToolName            string
	ToolHandlerFunction interface{}
	Description         string
	InputSchema         *jsonschema.Schema
	InputTypeName       string
	// for a tool to be available from a proxy, we need to set the ToolProxyId
	ToolProxyId string
}

type SdkToolProvider struct {
	toolName         string
	isDisabled       bool
	toolInitFunction interface{}
	contextType      reflect.Type
	contextTypeName  string
	toolDefinitions  []*SdkToolDefinition
	// the tool context retrieve from the tool init function
	// proxy id for proxy tool provider
	proxyId string
}

func newProxyToolProvider(proxyId string, proxyName string) (*SdkToolProvider, error) {
	toolProvider := &SdkToolProvider{
		toolName:         proxyName,
		isDisabled:       false,
		toolInitFunction: nil,
		contextType:      nil,
		contextTypeName:  "",
		toolDefinitions:  []*SdkToolDefinition{},
		proxyId:          proxyId,
	}
	return toolProvider, nil
}

// TODO:XXX: delete this
func (tp *SdkToolProvider) AddProxyTool(toolName string, description string, inputSchema interface{}) error {
	// Convert the interface{} to *jsonschema.Schema
	var schema *jsonschema.Schema
	switch s := inputSchema.(type) {
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

	// we need to check if the tool name is already registered
	for _, tool := range tp.toolDefinitions {
		if tool.ToolName == toolName {
			// we need to update the tool definition
			tool.Description = description
			tool.InputSchema = schema
			tool.ToolProxyId = tp.proxyId
			return nil
		}
	}

	// we create a new tool definition
	tp.toolDefinitions = append(tp.toolDefinitions, &SdkToolDefinition{
		ToolName:    toolName,
		Description: description,
		ToolProxyId: tp.proxyId,
		InputSchema: schema,
	})
	return nil
}

func mapToStruct(input map[string]interface{}, output interface{}) error {
	jsonBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonBytes, output)
}
