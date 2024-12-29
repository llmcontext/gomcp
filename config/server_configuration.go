package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/llmcontext/gomcp/jsonschema"
)

type ServerConfiguration struct {
	ConfigVersion int            `json:"v"`
	ServerInfo    ServerInfo     `json:"serverInfo"`
	Logging       *LoggingInfo   `json:"logging,omitempty"`
	Inspector     *InspectorInfo `json:"inspector,omitempty"`
	Prompts       *PromptConfig  `json:"prompts,omitempty"`
}

func LoadServerConfig(configFilePath string) (*ServerConfiguration, error) {

	// Check if the file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("server configuration file does not exist: %s", configFilePath)
	}

	// let's generate the schema from the config struct
	configSchema, err := jsonschema.GetSchemaFromAny(&ServerConfiguration{})
	if err != nil {
		return nil, fmt.Errorf("failed to generate schema from config struct: %v", err)
	}
	// let's check that the file is a valid json file
	jsonBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	err = jsonschema.ValidateJsonSchemaWithBytes(configSchema, jsonBytes)
	if err != nil {
		return nil, err
	}

	var config ServerConfiguration
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, err
	}

	// update the file path to be absolute
	if config.Logging != nil {
		config.Logging.UpdateFilePaths()
	}
	if config.Inspector != nil {
		config.Inspector.UpdateFilePaths()
	}

	return &config, nil
}
