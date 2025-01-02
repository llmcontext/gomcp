package config

import (
	"path/filepath"
)

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type LoggingInfo struct {
	File       string `json:"file,omitempty"`
	Level      string `json:"level,omitempty"`
	WithStderr bool   `json:"withStderr,omitempty"`
}

type InspectorInfo struct {
	Enabled           bool   `json:"enabled"`
	ListenAddress     string `json:"listenAddress"`
	ProtocolDebugFile string `json:"protocolDebugFile,omitempty"`
}

func updateFilePath(path string) string {
	if path == "" {
		return path
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(DefaultHubConfigurationDirectory, path)
	}
	return path
}

func (c *LoggingInfo) UpdateFilePaths() {
	c.File = updateFilePath(c.File)
}

func (c *InspectorInfo) UpdateFilePaths() {
	c.ProtocolDebugFile = updateFilePath(c.ProtocolDebugFile)
}
