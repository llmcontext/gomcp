package mcp

import (
	"github.com/llmcontext/gomcp/jsonrpc"
	"github.com/llmcontext/gomcp/protocol"
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
	result, err := protocol.CheckIsObject(response.Result, "result")
	if err != nil {
		return nil, err
	}

	// read protocol version
	protocolVersion, err := protocol.GetStringField(result, "protocolVersion")
	if err != nil {
		return nil, err
	}
	resp.ProtocolVersion = protocolVersion

	// read server info
	serverInfo, err := protocol.GetObjectField(result, "serverInfo")
	if err != nil {
		return nil, err
	}
	// read name
	name, err := protocol.GetStringField(serverInfo, "name")
	if err != nil {
		return nil, err
	}
	resp.ServerInfo.Name = name

	// read version
	version, err := protocol.GetStringField(serverInfo, "version")
	if err != nil {
		return nil, err
	}
	resp.ServerInfo.Version = version

	// read capabilities
	capabilities, err := protocol.CheckIsObject(result, "capabilities")
	if err != nil {
		return nil, err
	}
	// check if logging capability is present
	if _, ok := capabilities["logging"]; ok {
		resp.Capabilities.Logging = &ServerCapabilitiesLogging{}
	}

	// check if resources capability is present
	capabilitiesResources := protocol.GetOptionalObjectField(capabilities, "resources")
	if capabilitiesResources != nil {
		resp.Capabilities.Resources = &ServerCapabilitiesResources{}
		// check if listChanged is present
		listChanged := protocol.GetOptionalBoolField(capabilitiesResources, "listChanged")
		if listChanged != nil {
			resp.Capabilities.Resources.ListChanged = listChanged
		}
		// check if subscribe is present
		subscribe := protocol.GetOptionalBoolField(capabilitiesResources, "subscribe")
		if subscribe != nil {
			resp.Capabilities.Resources.Subscribe = subscribe
		}
	} else {
		resp.Capabilities.Resources = nil
	}

	// check if tools capability is present
	capabilitiesTools := protocol.GetOptionalObjectField(capabilities, "tools")
	if capabilitiesTools != nil {
		resp.Capabilities.Tools = &ServerCapabilitiesTools{}
		// check if listChanged is present
		listChanged := protocol.GetOptionalBoolField(capabilitiesTools, "listChanged")
		if listChanged != nil {
			resp.Capabilities.Tools.ListChanged = listChanged
		}
	} else {
		resp.Capabilities.Tools = nil
	}

	// check if prompts capability is present
	capabilitiesPrompts := protocol.GetOptionalObjectField(capabilities, "prompts")
	if capabilitiesPrompts != nil {
		resp.Capabilities.Prompts = &ServerCapabilitiesPrompts{}
		// check if listChanged is present
		listChanged := protocol.GetOptionalBoolField(capabilitiesPrompts, "listChanged")
		if listChanged != nil {
			resp.Capabilities.Prompts.ListChanged = listChanged
		}
	} else {
		resp.Capabilities.Prompts = nil
	}

	// check if logging capability is present
	logging := protocol.GetOptionalObjectField(capabilities, "logging")
	if logging != nil {
		resp.Capabilities.Logging = &ServerCapabilitiesLogging{}
	} else {
		resp.Capabilities.Logging = nil
	}

	return &resp, nil
}
