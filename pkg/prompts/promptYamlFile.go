package prompts

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/llmcontext/gomcp/pkg/jsonschema"
	"gopkg.in/yaml.v3"
)

// LoadPromptYamlFile reads and parses a yaml file containing a list of prompts
func LoadPromptYamlFile(filepath string) (*PromptList, error) {
	yamlData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// retrieve the schema for the PromptConfig struct
	configSchema, _, err := jsonschema.GetFullSchemaFromInterface(reflect.TypeOf(&PromptList{}))
	if err != nil {
		return nil, fmt.Errorf("error generating schema for PromptList")
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
	var config PromptList
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
