package messages

type JsonRpcResponseInitialize struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	Tools   ServerCapabilitiesTools   `json:"tools"`
	Prompts ServerCapabilitiesPrompts `json:"prompts"`
}

type ServerCapabilitiesTools struct {
	ListChanged bool `json:"listChanged"`
}

type ServerCapabilitiesPrompts struct {
	ListChanged bool `json:"listChanged"`
}
