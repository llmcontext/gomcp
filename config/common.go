package config

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type LoggingInfo struct {
	File              string `json:"file,omitempty"`
	Level             string `json:"level,omitempty"`
	WithStderr        bool   `json:"withStderr,omitempty"`
	ProtocolDebugFile string `json:"protocolDebugFile,omitempty"`
}

type InspectorInfo struct {
	Enabled       bool   `json:"enabled"`
	ListenAddress string `json:"listenAddress"`
}

type PromptConfig struct {
	File string `json:"file"`
}

type ToolConfig struct {
	Name          string      `json:"name"`
	IsDisabled    bool        `json:"isDisabled,omitempty"`
	Description   string      `json:"description,omitempty"`
	Configuration interface{} `json:"configuration"`
}
