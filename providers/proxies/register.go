package proxies

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/llmcontext/gomcp/jsonschema"
	"github.com/llmcontext/gomcp/registry"
)

func RegisterProxyServers(
	proxiesDirectory string,
	mcpServerRegistry *registry.McpServerRegistry) error {
	if _, err := os.Stat(proxiesDirectory); os.IsNotExist(err) {
		err := os.MkdirAll(proxiesDirectory, 0755)
		if err != nil {
			return fmt.Errorf("failed to create proxies directory: %v", err)
		}
		// no proxies to register
		return nil
	}

	// load all the proxy definitions
	files, err := os.ReadDir(proxiesDirectory)
	if err != nil {
		return err
	}

	// let's generate the schema from the config struct
	proxySchema, err := jsonschema.GetSchemaFromAny(&ProxyDefinition{})
	if err != nil {
		return err
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
		proxyPath := filepath.Join(proxiesDirectory, file.Name())

		// we unmarshal the file into a ProxyDefinition
		jsonBytes, err := os.ReadFile(proxyPath)
		if err != nil {
			return err
		}
		err = jsonschema.ValidateJsonSchemaWithBytes(proxySchema, jsonBytes)
		if err != nil {
			return err
		}

		var def ProxyDefinition
		err = json.Unmarshal(jsonBytes, &def)
		if err != nil {
			return err
		}

		// we register the proxy server
		// TODO: manage a version for the proxy server
		proxyHandler := &registry.McpServerLifecycle{
			Init: nil,
			End:  nil,
		}
		mcpServer, err := mcpServerRegistry.RegisterServer(def.ProxyName, "0.1", proxyHandler)
		if err != nil {
			return err
		}

		// register the proxy tools
		for _, tool := range def.Tools {
			proxyTool := NewProxyTool(def.ProxyId, &tool)

			proxyTool.register(mcpServer)
		}
	}
	return nil
}
