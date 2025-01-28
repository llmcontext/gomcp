package prompts

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/llmcontext/gomcp/pkg/jsonschema"
	"gopkg.in/yaml.v3"
)

// PromptConfig represents the root YAML structure
type PromptConfig struct {
	Prompts []*PromptDefinition `json:"prompts" yaml:"prompts"`
}

type PromptDefinition struct {
	Name        string               `json:"name" yaml:"name"`
	Description string               `json:"description" yaml:"description"`
	Arguments   []ArgumentDefinition `json:"arguments,omitempty" yaml:"arguments,omitempty"`
	Prompt      string               `json:"prompt" yaml:"prompt"`
}

type ArgumentDefinition struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`
}

// loadPrompts reads and parses the prompts.yaml file
func loadPrompts(filepath string) (*PromptConfig, error) {
	yamlData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// retrieve the schema for the PromptConfig struct
	configSchema, _, err := jsonschema.GetFullSchemaFromInterface(reflect.TypeOf(&PromptConfig{}))
	if err != nil {
		return nil, fmt.Errorf("error generating schema for toolInitFunction argument")
	}

	// unmarshal the yaml data into an interface{}
	var data interface{}
	if err := yaml.Unmarshal(yamlData, &data); err != nil {
		return nil, err
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("error converting to JSON: %w", err)
	}

	// validate the json data against the schema
	err = jsonschema.ValidateJsonSchemaWithBytes(configSchema, jsonData)
	if err != nil {
		return nil, err
	}

	// finally, unmarshal the json data back into the PromptConfig struct
	var config PromptConfig
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
