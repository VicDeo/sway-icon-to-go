package config

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
)

// ConfigDirProvider is an interface that provides the config directory.
type ConfigDirProvider interface {
	// GetConfigDir gets the config directory.
	GetConfigDir() string
}

// ConfigDirResolver is a struct that resolves the config directory.
type ConfigDirResolver struct{}

// NewConfigDirResolver creates a new ConfigDirResolver instance.
func NewConfigDirResolver() *ConfigDirResolver {
	return &ConfigDirResolver{}
}

// GetConfigDir gets the config directory.
func (c *ConfigDirResolver) GetConfigDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" && isValidDirectory(os.Getenv("XDG_CONFIG_HOME")) {
		return os.Getenv("XDG_CONFIG_HOME")
	}
	usr, err := user.Current()
	if err != nil {
		slog.Error("Error while getting current user", "error", err)
		return ""
	}
	home := usr.HomeDir
	return filepath.Join(home, ".config")
}

// ConfigFileResolver is a struct that resolves the config file path.
type ConfigFileResolver struct {
	ConfigDirProvider ConfigDirProvider
	Subdirs           []string
	Filename          string
}

// NewConfigFileResolver creates a new ConfigResolver instance.
func NewConfigFileResolver(provider ConfigDirProvider, subdirs []string, filename string) *ConfigFileResolver {
	return &ConfigFileResolver{
		ConfigDirProvider: provider,
		Subdirs:           subdirs,
		Filename:          filename,
	}
}

// Resolve resolves the config file path.
func (c *ConfigFileResolver) Resolve() (string, error) {
	for _, subdir := range c.Subdirs {
		fullPath := filepath.Join(c.ConfigDirProvider.GetConfigDir(), subdir)
		if !isValidDirectory(fullPath) {
			slog.Debug("Config directory does not exist or is not a directory", "path", fullPath)
			continue
		}

		configPath := filepath.Join(fullPath, c.Filename)
		if isValidFile(configPath) {
			return configPath, nil
		}
		slog.Debug("Config file not found", "path", configPath)
	}
	return "", fmt.Errorf("config file %s not found in any of the directories %v", c.Filename, c.Subdirs)
}

// isValidFile checks if the path is a file.
func isValidFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

// isValidDirectory checks if the path is a directory.
func isValidDirectory(path string) bool {
	dirInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return dirInfo.IsDir()
}
