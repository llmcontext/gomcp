package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/llmcontext/gomcp/defaults"
	"github.com/llmcontext/gomcp/jsonschema"
)

// configuration for the proxy
type ProxyConfiguration struct {
	ConfigurationFilePath string   `json:"-"`
	ConfigVersion         int      `json:"v"`
	WhatIsThat            string   `json:"what_is_that"`
	MoreInformation       string   `json:"more_information"`
	ProxyId               string   `json:"proxy_id"`
	ProgramName           string   `json:"program_name"`
	ProgramArgs           []string `json:"program_args"`
	LastStarted           string   `json:"last_started"`
}

func getDefaultProxyConfigurationPath(localDirectory string) string {
	return filepath.Join(localDirectory, defaults.DefaultProxyConfigPath)
}

func LoadProxyConfiguration(localDirectory string) (*ProxyConfiguration, error) {
	// Check if the file exists
	configPath := getDefaultProxyConfigurationPath(localDirectory)
	var proxyConfig = &ProxyConfiguration{
		ConfigurationFilePath: configPath,
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return proxyConfig, fmt.Errorf("proxy configuration file does not exist: %s", configPath)
	}

	// let's generate the schema from the config struct
	configSchema, err := jsonschema.GetSchemaFromAny(&ProxyConfiguration{})
	if err != nil {
		return proxyConfig, fmt.Errorf("failed to generate schema from config struct: %v", err)
	}
	// let's check that the file is a valid json file
	jsonBytes, err := os.ReadFile(configPath)
	if err != nil {
		return proxyConfig, err
	}

	err = jsonschema.ValidateJsonSchemaWithBytes(configSchema, jsonBytes)
	if err != nil {
		return proxyConfig, err
	}

	err = json.Unmarshal(jsonBytes, &proxyConfig)
	// we set the local directory here so that it is available to the config
	proxyConfig.ConfigurationFilePath = configPath
	if err != nil {
		return proxyConfig, err
	}

	return proxyConfig, nil
}

// SaveProxyConfig saves the proxy configuration to the file
func SaveProxyConfiguration(proxyConfig *ProxyConfiguration) error {
	configPath := proxyConfig.ConfigurationFilePath

	proxyConfig.LastStarted = time.Now().Format(time.RFC3339)
	json, err := json.MarshalIndent(proxyConfig, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, json, 0644)
}
