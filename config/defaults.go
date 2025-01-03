package config

import (
	"os"
	"path/filepath"
)

const (
	DefaultApplicationName = "gomcp"
	DefaultWsPort          = 8080
)

var DefaultHubConfigurationDirectory = filepath.Join(os.Getenv("HOME"), ".gomcp")
