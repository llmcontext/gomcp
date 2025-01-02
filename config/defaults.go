package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultApplicationName = "gomcp"
	DefaultMultiplexerPort = 8090
	DefaultWsPort          = 8080
	DefaultProxyDirectory  = "proxies"
)

var DefaultHubConfigurationDirectory = filepath.Join(os.Getenv("HOME"), ".gomcp")
