package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/llmcontext/gomcp/defaults"
)

type ProxyToolsRegistry struct {
	baseDirectory string
}

func NewProxyToolsRegistry() *ProxyToolsRegistry {
	var baseDirectory = filepath.Join(defaults.DefaultHubConfigurationDirectory, "proxy_tools")
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		os.MkdirAll(baseDirectory, 0755)
	}
	return &ProxyToolsRegistry{
		baseDirectory: baseDirectory,
	}
}

// called by the proxy to register a new proxy definition
func (t *ProxyToolsRegistry) AddProxyDefinition(def *ProxyDefinition) error {
	toolPath := filepath.Join(t.baseDirectory, fmt.Sprintf("%s.json", def.ProxyId))

	json, err := json.MarshalIndent(def, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(toolPath, json, 0644)
}
