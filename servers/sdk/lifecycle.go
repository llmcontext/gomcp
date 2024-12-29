package sdk

import (
	"context"

	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/jsonschema"
	"github.com/llmcontext/gomcp/registry"
	"github.com/llmcontext/gomcp/types"
)

func (s *SdkServerDefinition) serverInitFunction(ctx context.Context, logger types.Logger) error {
	var result interface{}
	var callErr, err error

	// check if we have a tool configuration data
	if s.toolConfigurationData != nil {
		result, callErr, err = callFunction(s.toolsInitFunction, ctx, s.toolConfigurationData)
	} else {
		result, callErr, err = callFunction(s.toolsInitFunction, ctx)
	}
	if err != nil {
		return err
	}
	if callErr != nil {
		return callErr
	}
	logger.Info("tool provider initialized", types.LogArg{
		"result": result,
	})

	// check if result as a property called logger of type types.Logger
	if logger, ok := result.(types.Logger); ok {
		logger.Info("tool provider initialized", types.LogArg{
			"result": result,
		})
	}

	// we store the tool context
	s.toolContext = result

	// we pass it to all the tools
	for _, tool := range s.toolDefinitions {
		tool.toolContext = s.toolContext
	}

	return nil
}

func (s *SdkServerDefinition) serverEndFunction(ctx context.Context, logger types.Logger) error {
	return nil
}

func (t *SdkToolDefinition) toolInitFunction(ctx context.Context, logger types.Logger) error {
	return nil
}

func (t *SdkToolDefinition) toolProcessFunction(
	ctx context.Context,
	toolArgs map[string]interface{},
	result types.ToolCallResult,
	logger types.Logger,
	errChan chan *jsonrpc.JsonRpcError,
) error {

	// let's check if the arguments match the schema
	err := jsonschema.ValidateJsonSchemaWithObject(t.inputSchema, toolArgs)
	if err != nil {
		return err
	}

	// create a new context with the logger
	goCtx := types.ContextWithLogger(ctx, logger)

	// let's create the output
	output := registry.NewToolCallResult()

	go func() {
		_, callErr, err := callFunction(t.toolHandlerFunction, goCtx, t.toolContext, toolArgs, output)
		if err != nil {
			errChan <- &jsonrpc.JsonRpcError{
				Code:    jsonrpc.RpcInternalError,
				Message: err.Error(),
			}
		} else if callErr != nil {
			errChan <- &jsonrpc.JsonRpcError{
				Code:    jsonrpc.RpcInternalError,
				Message: callErr.Error(),
			}
		} else {
			errChan <- nil
		}
	}()

	return nil
}

func (t *SdkToolDefinition) toolEndFunction(ctx context.Context, logger types.Logger) error {
	return nil
}
