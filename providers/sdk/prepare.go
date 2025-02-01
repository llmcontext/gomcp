package sdk

import (
	"context"
	"fmt"
	"reflect"

	"github.com/llmcontext/gomcp/pkg/jsonschema"
	"github.com/llmcontext/gomcp/types"
)

// stitch everything together
// so that we can use the server and tools
func (s *SdkServerDefinition) Prepare() error {
	// we setup the server
	// check that the tools are valid
	err := s.setupServer()
	if err != nil {
		return fmt.Errorf("failed to setup MCP server: %v", err)
	}

	// we add all the tools to the tools registry
	for _, tool := range s.toolDefinitions {
		err := tool.setupTool(s)
		if err != nil {
			return fmt.Errorf("failed to setup tool %s: %v", tool.ToolName, err)
		}
	}

	return nil
}

// TODO: check that the configuration is valid eg. s.toolsInitFunction exists etc...
func (s *SdkServerDefinition) setupServer() error {
	// get the type of the configuration
	configurationType := reflect.TypeOf(s.toolConfigurationData)

	// Validate that toolHandler is a function
	fnType := reflect.TypeOf(s.toolsInitFunction)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("toolInitFunction must be a function")
	}

	// the function must have 1 or 2 arguments: context and optional config
	if fnType.NumIn() != 1 && fnType.NumIn() != 2 {
		return fmt.Errorf("toolInitFunction must have 1 or 2 arguments")
	}

	// the first argument must be a golang context
	goContextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if fnType.In(0) != goContextType {
		return fmt.Errorf("first argument must be context.Context")
	}

	var configType reflect.Type = nil

	// if 1 argument is provided, it must be a pointer to a struct
	if fnType.NumIn() == 2 {
		if fnType.In(1).Kind() != reflect.Ptr {
			return fmt.Errorf("toolInitFunction argument must be a pointer to a struct")
		}

		// get the type of the configuration the pointer is pointing to (.Elem())
		configType = fnType.In(1)

		// check if the type is the same as the configuration type
		if configType != configurationType {
			return fmt.Errorf("toolInitFunction argument must be a pointer to a struct of type %s, but got %s", configurationType.String(), configType.String())
		}
	}

	// the function must return a tool context, error
	if fnType.NumOut() != 2 || fnType.Out(0).Kind() != reflect.Ptr || fnType.Out(1).Kind() != reflect.Interface {
		return fmt.Errorf("toolInitFunction must return a context, error")
	}

	// check that the second output is an error
	if fnType.Out(1).String() != "error" {
		return fmt.Errorf("toolInitFunction second return value must be an error")
	}

	// the first return value must be a context
	if fnType.Out(0).Kind() != reflect.Ptr {
		return fmt.Errorf("toolInitFunction first return value must be a pointer to a context")
	}

	// the context must be a struct
	if fnType.Out(0).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("toolInitFunction first return value must be a pointer to a struct")
	}
	returnedContextType := fnType.Out(0).Elem()
	returnedContextTypeName := returnedContextType.Name()
	s.contextType = returnedContextType
	s.contextTypeName = returnedContextTypeName

	return nil
}

func (tool *SdkToolDefinition) setupTool(serverDefinition *SdkServerDefinition) error {
	// Validate that toolHandler is a function
	fnType := reflect.TypeOf(tool.toolHandlerFunction)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("toolHandler must be a function")
	}

	// the function must have 4 arguments:
	// the golang context
	// the tool context
	// the input
	// the output
	if fnType.NumIn() != 4 {
		return fmt.Errorf("toolHandler for %s must have 4 arguments", tool.ToolName)
	}

	// the first argument must be a golang context
	goContextType := reflect.TypeOf((*context.Context)(nil)).Elem()
	if fnType.In(0) != goContextType {
		return fmt.Errorf("toolHandler for %s first argument must be a golang context", tool.ToolName)
	}

	// the second argument must be a pointer to the tool context type
	if fnType.In(1).Kind() != reflect.Ptr || fnType.In(1).Elem() != serverDefinition.contextType {
		return fmt.Errorf("toolHandler for %s second argument must be a pointer to the context type: %s", tool.ToolName, serverDefinition.contextTypeName)
	}

	// the third argument must be a pointer to a struct
	if fnType.In(2).Kind() != reflect.Ptr || fnType.In(2).Elem().Kind() != reflect.Struct {
		return fmt.Errorf("toolHandler for %s third argument must be a pointer to a struct", tool.ToolName)
	}
	// we need to get the schema of the third argument
	inputSchema, inputTypeName, err := jsonschema.GetSchemaFromType(fnType.In(2))
	if err != nil {
		return fmt.Errorf("error generating schema for toolHandler for %s third argument", tool.ToolName)
	}

	// the fourth argument must be an implementation of types.ToolCallResult
	toolCallResultType := reflect.TypeOf((*types.ToolCallResult)(nil)).Elem()
	if !fnType.In(3).Implements(toolCallResultType) {
		return fmt.Errorf("toolHandler for %s fourth argument must implement types.ToolCallResult but is %s", tool.ToolName, fnType.In(3).String())
	}

	// the function must return an error
	if fnType.NumOut() != 1 || fnType.Out(0).String() != "error" {
		return fmt.Errorf("toolHandler for %s must return an error", tool.ToolName)
	}

	// Store the function for later use
	tool.InputSchema = inputSchema
	tool.inputTypeName = inputTypeName

	return nil
}
