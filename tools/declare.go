package tools

import (
	"context"
	"fmt"
	"reflect"

	"github.com/invopop/jsonschema"
	"github.com/llmcontext/gomcp/logger"
)

type ToolDefinition struct {
	ToolName            string
	ToolHandlerFunction interface{}
	Description         string
	InputSchema         *jsonschema.Schema
	InputTypeName       string
}

type ToolProvider struct {
	toolName         string
	isDisabled       bool
	configSchema     *jsonschema.Schema
	configTypeName   string
	configType       reflect.Type
	toolInitFunction interface{}
	contextType      reflect.Type
	contextTypeName  string
	toolDefinitions  []ToolDefinition
	// the tool context retrieve from the tool init function
	toolContext interface{}
}

func getSchemaFromType(t reflect.Type) (*jsonschema.Schema, string, error) {
	var typeName = t.Elem().Name()
	if typeName == "" {
		return nil, "", fmt.Errorf("type name is empty")
	}

	reflector := jsonschema.Reflector{}
	schema := reflector.ReflectFromType(t)
	if schema == nil {
		return nil, "", fmt.Errorf("error generating schema")
	}

	schemaType := schema.Definitions[typeName]
	if schemaType == nil {
		return nil, "", fmt.Errorf("no schema for definition found")
	}
	logger.Info("schema", logger.Arg{
		"type":   typeName,
		"schema": schemaType,
	})

	return schemaType, typeName, nil
}

func DeclareToolProvider(toolName string, toolInitFunction interface{}) (*ToolProvider, error) {
	// we initialize the tool provider with nil values
	toolProvider := &ToolProvider{
		toolName:         toolName,
		isDisabled:       false,
		configSchema:     nil,
		configTypeName:   "",
		configType:       nil,
		toolInitFunction: toolInitFunction,
		contextType:      nil,
		contextTypeName:  "",
		toolDefinitions:  []ToolDefinition{},
	}

	// Validate that toolHandler is a function
	fnType := reflect.TypeOf(toolInitFunction)
	if fnType.Kind() != reflect.Func {
		return nil, fmt.Errorf("toolInitFunctiom must be a function")
	}

	// the function must have 1 or 2 arguments: context and optional config
	if fnType.NumIn() != 1 && fnType.NumIn() != 2 {
		return nil, fmt.Errorf("toolInitFunctiom must have 1 or 2 arguments")
	}

	// the first argument must be a golang context
	goContextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if fnType.In(0) != goContextType {
		return nil, fmt.Errorf("first argument must be context.Context")
	}

	var configType reflect.Type = nil

	// if 1 argument is provided, it must be a pointer to a struct
	if fnType.NumIn() == 2 {
		if fnType.In(1).Kind() != reflect.Ptr {
			return nil, fmt.Errorf("toolInitFunctiom argument must be a pointer to a struct")
		}
		configType = fnType.In(1)
		configSchema, configTypeName, err := getSchemaFromType(configType)
		if err != nil {
			return nil, fmt.Errorf("error generating schema for toolInitFunctiom argument")
		}
		// we store the config schema, type name and type
		toolProvider.configSchema = configSchema
		toolProvider.configTypeName = configTypeName
		toolProvider.configType = configType
	}

	// the function must return a tool context, error
	if fnType.NumOut() != 2 || fnType.Out(0).Kind() != reflect.Ptr || fnType.Out(1).Kind() != reflect.Interface {
		return nil, fmt.Errorf("toolInitFunctiom must return a context, error")
	}

	// check that the second output is an error
	if fnType.Out(1).String() != "error" {
		return nil, fmt.Errorf("toolInitFunctiom second return value must be an error")
	}

	// the first return value must be a context
	if fnType.Out(0).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("toolInitFunctiom first return value must be a pointer to a context")
	}

	// the context must be a struct
	if fnType.Out(0).Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("toolInitFunctiom first return value must be a pointer to a struct")
	}
	returnedContextType := fnType.Out(0).Elem()
	returnedContextTypeName := returnedContextType.Name()
	toolProvider.contextType = returnedContextType
	toolProvider.contextTypeName = returnedContextTypeName

	return toolProvider, nil
}

func (tp *ToolProvider) AddTool(toolName string, description string, toolHandler interface{}) error {
	// Validate that toolHandler is a function
	fnType := reflect.TypeOf(toolHandler)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("toolHandler must be a function")
	}

	// the function must have 4 arguments:
	// the golang context
	// the tool context
	// the input
	// the output
	if fnType.NumIn() != 4 {
		return fmt.Errorf("toolHandler for %s must have 4 arguments", toolName)
	}

	// the first argument must be a golang context
	goContextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if fnType.In(0) != goContextType {
		return fmt.Errorf("toolHandler for %s first argument must be a golang context", toolName)
	}

	// the second argument must be a pointer to the tool context type
	if fnType.In(1).Kind() != reflect.Ptr || fnType.In(1).Elem() != tp.contextType {
		return fmt.Errorf("toolHandler for %s second argument must be a pointer to the context type: %s", toolName, tp.contextTypeName)
	}

	// the third argument must be a pointer to a struct
	if fnType.In(2).Kind() != reflect.Ptr || fnType.In(2).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("toolHandler for %s third argument must be a pointer to a struct", toolName)
	}
	// we need to get the schema of the third argument
	inputSchema, inputTypeName, err := getSchemaFromType(fnType.In(2))
	if err != nil {
		return fmt.Errorf("error generating schema for toolHandler for %s third argument", toolName)
	}

	// the fourth argument must be a pointer to tools.ToolCallResult
	toolCallResultType := reflect.TypeOf((*ToolCallResult)(nil))
	if fnType.In(3) != toolCallResultType {
		return fmt.Errorf("toolHandler for %s fourth argument must be a *tools.ToolCallResult", toolName)
	}

	// the function must return an error
	if fnType.NumOut() != 1 || fnType.Out(0).String() != "error" {
		return fmt.Errorf("toolHandler for %s must return an error", toolName)
	}

	// Store the function for later use
	tp.toolDefinitions = append(tp.toolDefinitions, ToolDefinition{
		ToolName:            toolName,
		Description:         description,
		ToolHandlerFunction: toolHandler,
		InputSchema:         inputSchema,
		InputTypeName:       inputTypeName,
	})
	return nil
}
