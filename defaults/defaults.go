package defaults

import (
	"os"
	"path/filepath"
)

const (
	DefaultApplicationName = "gomcp"
	DefaultMultiplexerPort = 8090
	DefaultWsPort          = 8080
	DefaultProxyConfigPath = "gomcp-proxy.json"
)

var DefaultHubConfigurationDirectory = filepath.Join(os.Getenv("HOME"), ".gomcp")
