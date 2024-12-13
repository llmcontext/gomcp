package mcp

import (
	"github.com/llmcontext/gomcp/jsonrpc"
)

type JsonRpcResponseInitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	Tools     *ServerCapabilitiesTools     `json:"tools,omitempty"`
	Prompts   *ServerCapabilitiesPrompts   `json:"prompts,omitempty"`
	Logging   *ServerCapabilitiesLogging   `json:"logging,omitempty"`
	Resources *ServerCapabilitiesResources `json:"resources,omitempty"`
}

type ServerCapabilitiesTools struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type ServerCapabilitiesPrompts struct {
	ListChanged *bool `json:"listChanged,omitempty"`
}

type ServerCapabilitiesLogging struct {
}

type ServerCapabilitiesResources struct {
	ListChanged *bool `json:"listChanged,omitempty"`
	Subscribe   *bool `json:"subscribe,omitempty"`
}

func ParseJsonRpcResponseInitialize(response *jsonrpc.JsonRpcResponse) (*JsonRpcResponseInitializeResult, error) {
	resp := JsonRpcResponseInitializeResult{}

	// parse params
	result, err := checkIsObject(response.Result, "result")
	if err != nil {
		return nil, err
	}

	// read protocol version
	protocolVersion, err := getStringField(result, "protocolVersion")
	if err != nil {
		return nil, err
	}
	resp.ProtocolVersion = protocolVersion

	// read server info
	serverInfo, err := getObjectField(result, "serverInfo")
	if err != nil {
		return nil, err
	}
	// read name
	name, err := getStringField(serverInfo, "name")
	if err != nil {
		return nil, err
	}
	resp.ServerInfo.Name = name

	// read version
	version, err := getStringField(serverInfo, "version")
	if err != nil {
		return nil, err
	}
	resp.ServerInfo.Version = version

	// read capabilities
	capabilities, err := checkIsObject(result, "capabilities")
	if err != nil {
		return nil, err
	}
	// check if logging capability is present
	if _, ok := capabilities["logging"]; ok {
		resp.Capabilities.Logging = &ServerCapabilitiesLogging{}
	}

	// check if resources capability is present
	capabilitiesResources := getOptionalObjectField(capabilities, "resources")
	if capabilitiesResources != nil {
		resp.Capabilities.Resources = &ServerCapabilitiesResources{}
		// check if listChanged is present
		listChanged := getOptionalBoolField(capabilitiesResources, "listChanged")
		if listChanged != nil {
			resp.Capabilities.Resources.ListChanged = listChanged
		}
		// check if subscribe is present
		subscribe := getOptionalBoolField(capabilitiesResources, "subscribe")
		if subscribe != nil {
			resp.Capabilities.Resources.Subscribe = subscribe
		}
	} else {
		resp.Capabilities.Resources = nil
	}

	// check if tools capability is present
	capabilitiesTools := getOptionalObjectField(capabilities, "tools")
	if capabilitiesTools != nil {
		resp.Capabilities.Tools = &ServerCapabilitiesTools{}
		// check if listChanged is present
		listChanged := getOptionalBoolField(capabilitiesTools, "listChanged")
		if listChanged != nil {
			resp.Capabilities.Tools.ListChanged = listChanged
		}
	} else {
		resp.Capabilities.Tools = nil
	}

	// check if prompts capability is present
	capabilitiesPrompts := getOptionalObjectField(capabilities, "prompts")
	if capabilitiesPrompts != nil {
		resp.Capabilities.Prompts = &ServerCapabilitiesPrompts{}
		// check if listChanged is present
		listChanged := getOptionalBoolField(capabilitiesPrompts, "listChanged")
		if listChanged != nil {
			resp.Capabilities.Prompts.ListChanged = listChanged
		}
	} else {
		resp.Capabilities.Prompts = nil
	}

	// check if logging capability is present
	logging := getOptionalObjectField(capabilities, "logging")
	if logging != nil {
		resp.Capabilities.Logging = &ServerCapabilitiesLogging{}
	} else {
		resp.Capabilities.Logging = nil
	}

	return &resp, nil
}
