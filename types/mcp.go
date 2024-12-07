package types

type ToolProvider interface {
	AddTool(toolName string, description string, toolHandler interface{}) error
}

type ToolRegistry interface {
	DeclareToolProvider(toolName string, toolInitFunction interface{}) (ToolProvider, error)
}

type ModelContextProtocol interface {
	StdioTransport() Transport
	GetToolRegistry() ToolRegistry
	Start(serverName string, serverVersion string, transport Transport) error
}
