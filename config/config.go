package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/llmcontext/gomcp/utils"
)

type Config struct {
	ServerInfo ServerInfo     `json:"serverInfo"`
	Logging    LoggingInfo    `json:"logging,omitempty"`
	Inspector  *InspectorInfo `json:"inspector,omitempty"`
	Tools      []ToolConfig   `json:"tools,omitempty"`
	Prompts    *PromptConfig  `json:"prompts,omitempty"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ToolConfig struct {
	Name          string      `json:"name"`
	IsDisabled    bool        `json:"isDisabled,omitempty"`
	Description   string      `json:"description,omitempty"`
	Configuration interface{} `json:"configuration"`
}

type LoggingInfo struct {
	File              string `json:"file,omitempty"`
	Level             string `json:"level,omitempty"`
	WithStderr        bool   `json:"withStderr,omitempty"`
	ProtocolDebugFile string `json:"protocolDebugFile,omitempty"`
}

type InspectorInfo struct {
	Enabled       bool   `json:"enabled"`
	ListenAddress string `json:"listenAddress"`
}

type PromptConfig struct {
	File string `json:"file"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(configFilePath string) (*Config, error) {

	// Check if the file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configFilePath)
	}

	// let's generate the schema from the config struct
	configSchema := jsonschema.Reflect(&Config{})
	if configSchema == nil {
		return nil, fmt.Errorf("failed to generate schema from config struct")
	}
	// let's check that the file is a valid json file
	jsonBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	err = utils.ValidateJsonSchemaWithBytes(configSchema, jsonBytes)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(jsonBytes, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
