package proxy

type ProxyDefinition struct {
	ProxyId            string                 `json:"proxyId"`
	WorkingDirectory   string                 `json:"workingDirectory"`
	ProxyServerName    string                 `json:"proxyServerName"`
	ProxyServerVersion string                 `json:"proxyServerVersion"`
	ProgramName        string                 `json:"programName"`
	ProgramArguments   []string               `json:"programArguments"`
	Tools              []*ProxyToolDefinition `json:"tools"`
}

type ProxyToolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}
