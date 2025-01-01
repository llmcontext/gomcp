package config

import (
	"path/filepath"

	"github.com/llmcontext/gomcp/defaults"
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

// TODO: delete and have prompt part of a preset server
type PromptConfig struct {
	File string `json:"file"`
}

func updateFilePath(path string) string {
	if path == "" {
		return path
	}
	if !filepath.IsAbs(path) {
		path = filepath.Join(defaults.DefaultHubConfigurationDirectory, path)
	}
	return path
}

func (c *LoggingInfo) UpdateFilePaths() {
	c.File = updateFilePath(c.File)
}

func (c *InspectorInfo) UpdateFilePaths() {
	c.ProtocolDebugFile = updateFilePath(c.ProtocolDebugFile)
}
