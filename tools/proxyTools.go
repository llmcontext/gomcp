package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"
	"github.com/llmcontext/gomcp/defaults"
	"github.com/llmcontext/gomcp/utils"
)

type ProxyTools struct {
	baseDirectory string
}

type ProxyDefinition struct {
	ProxyId          string                `json:"proxyId"`
	WorkingDirectory string                `json:"workingDirectory"`
	ProxyName        string                `json:"proxyName"`
	ProgramName      string                `json:"programName"`
	ProgramArguments []string              `json:"programArguments"`
	Tools            []ProxyToolDefinition `json:"tools"`
}

type ProxyToolDefinition struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

func NewProxyTools() *ProxyTools {
	var baseDirectory = filepath.Join(defaults.DefaultHubConfigurationDirectory, defaults.DefaultProxyToolsDirectory)
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		os.MkdirAll(baseDirectory, 0755)
	}
	return &ProxyTools{
		baseDirectory: baseDirectory,
	}
}

func (t *ProxyTools) RegisterProxyTools(toolsRegistry *ToolsRegistry) error {
	// load all the proxy definitions
	files, err := os.ReadDir(t.baseDirectory)
	if err != nil {
		return err
	}

	// let's generate the schema from the config struct
	proxySchema := jsonschema.Reflect(&ProxyDefinition{})
	if proxySchema == nil {
		return fmt.Errorf("failed to generate schema from config struct")
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

		proxyPath := filepath.Join(t.baseDirectory, file.Name())

		// we unmarshal the file into a ProxyDefinition
		jsonBytes, err := os.ReadFile(proxyPath)
		if err != nil {
			return err
		}
		err = utils.ValidateJsonSchemaWithBytes(proxySchema, jsonBytes)
		if err != nil {
			return err
		}

		var def ProxyDefinition
		err = json.Unmarshal(jsonBytes, &def)
		if err != nil {
			return err
		}

		toolProvider, err := toolsRegistry.RegisterProxyToolProvider(def.ProxyId, def.ProxyName)
		if err != nil {
			return err
		}

		// register the proxy tools
		for _, tool := range def.Tools {
			err := toolProvider.AddProxyTool(tool.Name, tool.Description, tool.InputSchema)
			if err != nil {
				return err
			}
		}
		// we need to prepare the tool provider so that it can be used by the hub
		err = toolsRegistry.PrepareProxyToolProvider(toolProvider)
		if err != nil {
			return err
		}
	}
	return nil
}
