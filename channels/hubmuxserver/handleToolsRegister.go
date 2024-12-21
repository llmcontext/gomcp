package hubmuxserver

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol/mux"
	"github.com/llmcontext/gomcp/types"
)

func (s *MuxSession) handleToolsRegister(request *jsonrpc.JsonRpcRequest) error {
	params, err := mux.ParseJsonRpcRequestToolsRegisterParams(request)
	if err != nil {
		return err
	}
	s.logger.Info("Tools register", types.LogArg{
		"tools": params.Tools,
	})
	toolProvider, err := s.toolsRegistry.RegisterProxyToolProvider(s.proxyId, s.proxyName)
	if err != nil {
		s.logger.Error("Failed to register proxy tool provider", types.LogArg{
			"error": err,
		})
		return err
	}
	for _, tool := range params.Tools {
		err := toolProvider.AddProxyTool(tool.Name, tool.Description, tool.InputSchema)
		if err != nil {
			s.logger.Error("Failed to add proxy tool", types.LogArg{
				"error": err,
			})
			return err
		}
	}

	// we need to prepare the tool provider so that it can be used by the hub
	err = s.toolsRegistry.PrepareProxyToolProvider(toolProvider)
	if err != nil {
		s.logger.Error("Failed to prepare proxy tool provider", types.LogArg{
			"error": err,
		})
		return err
	}

	// s.events.EventNewProxyTools()

	return nil
}
