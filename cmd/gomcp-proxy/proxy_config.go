package main

import (
	"encoding/json"
	"os"
)

// configuration for the proxy
type ProxyConfig struct {
	WhatIsThat      string   `json:"what_is_that"`
	MoreInformation string   `json:"more_information"`
	ProxyId         string   `json:"proxy_id"`
	MuxAddress      string   `json:"mux_address"`
	ProgramName     string   `json:"program_name"`
	ProgramArgs     []string `json:"program_args"`
}

// SaveProxyConfig saves the proxy configuration to the file
func SaveProxyConfig(configPath string, proxyConfig *ProxyConfig) error {
	json, err := json.MarshalIndent(proxyConfig, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, json, 0644)
}

func LoadProxyConfig(configPath string) (*ProxyConfig, error) {
	config := &ProxyConfig{}

	// check if the config file exists
	// if it does not exist, return nil
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, nil
	}

	// read the config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(content, config)
	return config, nil
}
