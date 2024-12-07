package types

type ToolProvider interface {
	AddTool(toolName string, description string, toolHandler interface{}) error
}

type ModelContextProtocol interface {
	StdioTransport() Transport
	DeclareToolProvider(toolName string, toolInitFunction interface{}) (ToolProvider, error)
	Start(serverName string, serverVersion string, transport Transport) error
}
