package proxies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/llmcontext/gomcp/config"
	"github.com/llmcontext/gomcp/jsonschema"
)

type ProxyRegistry struct {
	baseDirectory string
	proxies       []*ProxyDefinition
}

func NewProxyRegistry() (*ProxyRegistry, error) {
	baseDirectory := filepath.Join(config.DefaultHubConfigurationDirectory, config.DefaultProxyDirectory)
	proxies, err := getListProxies(baseDirectory)
	if err != nil {
		return nil, err
	}
	return &ProxyRegistry{
		baseDirectory: baseDirectory,
		proxies:       proxies,
	}, nil
}

func (r *ProxyRegistry) GetProxies() []*ProxyDefinition {
	return r.proxies
}

func (r *ProxyRegistry) GetProxy(proxyId string) *ProxyDefinition {
	for _, proxy := range r.proxies {
		if proxy.ProxyId == proxyId {
			return proxy
		}
	}
	return nil
}

func (r *ProxyRegistry) GetTool(toolName string) *ProxyToolDefinition {
	for _, proxy := range r.proxies {
		for _, tool := range proxy.Tools {
			if tool.Name == toolName {
				return tool
			}
		}
	}
	return nil
}

func (r *ProxyRegistry) AddProxy(proxy *ProxyDefinition) error {
	// we need to check if the proxy already exists
	for _, p := range r.proxies {
		if p.ProxyId == proxy.ProxyId {
			return fmt.Errorf("proxy with same id already exists")
		}
	}
	// check if we have existing proxy with same tool name
	for _, p := range r.proxies {
		for _, tool := range p.Tools {
			for _, newTool := range proxy.Tools {
				if tool.Name == newTool.Name {
					return fmt.Errorf("proxy with same tool name already exists (%s)", newTool.Name)
				}
			}
		}
	}

	// save the proxy to the file
	proxyPath := filepath.Join(r.baseDirectory, proxy.ProxyId+".json")

	// format json to be readable
	jsonBytes, err := json.MarshalIndent(proxy, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(proxyPath, jsonBytes, 0644)
	if err != nil {
		return err
	}

	// we add the proxy to the registry
	r.proxies = append(r.proxies, proxy)
	return nil
}

func getListProxies(baseDirectory string) ([]*ProxyDefinition, error) {
	var proxies = []*ProxyDefinition{}
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(baseDirectory, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to create proxies directory: %v", err)
		}
		// no proxies to register
		return proxies, nil
	}

	// load all the proxy definitions
	files, err := os.ReadDir(baseDirectory)
	if err != nil {
		return nil, err
	}

	// let's generate the schema from the config struct
	proxySchema, err := jsonschema.GetSchemaFromAny(&ProxyDefinition{})
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		// skip directories
		if file.IsDir() {
			continue
		}

		// ensure the file is a json file
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		// description file for the proxy
		proxyPath := filepath.Join(baseDirectory, file.Name())

		// we unmarshal the file into a ProxyDefinition
		jsonBytes, err := os.ReadFile(proxyPath)
		if err != nil {
			return nil, err
		}
		err = jsonschema.ValidateJsonSchemaWithBytes(proxySchema, jsonBytes)
		if err != nil {
			return nil, err
		}

		var def ProxyDefinition
		err = json.Unmarshal(jsonBytes, &def)
		if err != nil {
			return nil, err
		}
		for _, tool := range def.Tools {
			err = setJsonSchema(tool)
			if err != nil {
				return nil, err
			}
		}
		proxies = append(proxies, &def)
	}
	return proxies, nil
}

func setJsonSchema(tool *ProxyToolDefinition) error {
	schema, err := jsonschema.ToJsonSchema(tool.InputSchema)
	if err != nil {
		return err
	}
	tool.JsonSchema = schema
	return nil
}

func (p *ProxyDefinition) GetTools() []*ProxyToolDefinition {
	return p.Tools
}

func (r *ProxyRegistry) Prepare() error {
	// TODO: to be implemented
	return nil
}
