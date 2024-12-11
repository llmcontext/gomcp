package messages

import (
	"fmt"

	"github.com/llmcontext/gomcp/jsonrpc"
)

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

func ParseJsonRpcResponseInitialize(response *jsonrpc.JsonRpcResponse) (*JsonRpcResponseInitialize, error) {
	resp := JsonRpcResponseInitialize{}

	// parse params
	if response.Result == nil {
		return nil, fmt.Errorf("missing result")
	}

	result, ok := response.Result.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("result must be an object")
	}

	// read protocol version
	protocolVersion, ok := result["protocolVersion"].(string)
	if !ok {
		return nil, fmt.Errorf("missing protocolVersion")
	}
	resp.ProtocolVersion = protocolVersion

	// read server info
	serverInfo, ok := result["serverInfo"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing serverInfo")
	}
	// read name
	name, ok := serverInfo["name"].(string)
	if !ok {
		return nil, fmt.Errorf("missing name")
	}
	resp.ServerInfo.Name = name

	// read version
	version, ok := serverInfo["version"].(string)
	if !ok {
		return nil, fmt.Errorf("missing version")
	}
	resp.ServerInfo.Version = version

	// read capabilities
	capabilities, ok := result["capabilities"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("missing or wrong capabilities")
	}
	// check if logging capability is present
	if _, ok := capabilities["logging"]; ok {
		resp.Capabilities.Logging = &ServerCapabilitiesLogging{}
	}

	// check if resources capability is present
	if capabilitiesResources, ok := capabilities["resources"]; ok {
		// check if resources is an object
		props, ok := capabilitiesResources.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("resources must be an object")
		}
		resp.Capabilities.Resources = &ServerCapabilitiesResources{}
		// check if listChanged is present
		if listChanged, ok := props["listChanged"].(bool); ok {
			resp.Capabilities.Resources.ListChanged = &listChanged
		}
		// check if subscribe is present
		if subscribe, ok := props["subscribe"].(bool); ok {
			resp.Capabilities.Resources.Subscribe = &subscribe
		}
	} else {
		resp.Capabilities.Resources = nil
	}

	// check if tools capability is present
	if capabilitiesTools, ok := capabilities["tools"]; ok {
		// check if tools is an object
		props, ok := capabilitiesTools.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("tools must be an object")
		}
		resp.Capabilities.Tools = &ServerCapabilitiesTools{}
		// check if listChanged is present
		if listChanged, ok := props["listChanged"].(bool); ok {
			resp.Capabilities.Tools.ListChanged = &listChanged
		}
	} else {
		resp.Capabilities.Tools = nil
	}

	// check if prompts capability is present
	if capabilitiesPrompts, ok := capabilities["prompts"]; ok {
		// check if prompts is an object
		props, ok := capabilitiesPrompts.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("prompts must be an object")
		}
		resp.Capabilities.Prompts = &ServerCapabilitiesPrompts{}

		// check if listChanged is present
		if listChanged, ok := props["listChanged"].(bool); ok {
			resp.Capabilities.Prompts.ListChanged = &listChanged
		}
	} else {
		resp.Capabilities.Prompts = nil
	}

	// check if logging capability is present
	if capabilitiesLogging, ok := capabilities["logging"]; ok {
		// check if logging is an object
		_, ok := capabilitiesLogging.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("logging must be an object")
		}
		resp.Capabilities.Logging = &ServerCapabilitiesLogging{}
	} else {
		resp.Capabilities.Logging = nil
	}

	return &resp, nil
}
