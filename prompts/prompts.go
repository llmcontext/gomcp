package prompts

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/llmcontext/gomcp/utils"
	"gopkg.in/yaml.v3"
)

// PromptConfig represents the root YAML structure
type PromptConfig struct {
	Prompts []struct {
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description" yaml:"description"`
		Arguments   []struct {
			Name        string `json:"name" yaml:"name"`
			Type        string `json:"type" yaml:"type"`
			Description string `json:"description" yaml:"description"`
			Required    bool   `json:"required" yaml:"required"`
		} `json:"arguments,omitempty" yaml:"arguments,omitempty"`
		Prompt string `json:"prompt" yaml:"prompt"`
	} `json:"prompts" yaml:"prompts"`
}

// LoadPrompts reads and parses the prompts.yaml file
func LoadPrompts(filepath string) (*PromptConfig, error) {
	yamlData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// retrieve the schema for the PromptConfig struct
	configSchema, _, err := utils.GetSchemaFromType(reflect.TypeOf(&PromptConfig{}))
	if err != nil {
		return nil, fmt.Errorf("error generating schema for toolInitFunctiom argument")
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
	err = utils.ValidateJsonSchemaWithBytes(configSchema, jsonData)
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
