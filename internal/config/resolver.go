package config

import (
	"fmt"
	"log/slog"
	"os"
	"os/user"
	"path/filepath"
)

// PathValidator is an interface that validates the path.
type PathValidator interface {
	IsValidDirectory(path string) bool
	IsValidFile(path string) bool
}

// pathValidator is a struct that validates the path.
type pathValidator struct{}

// NewPathValidator creates a new PathValidator instance.
func NewPathValidator() *pathValidator {
	return &pathValidator{}
}

// IsValidFile checks if the path is a file.
func (v *pathValidator) IsValidFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !fileInfo.IsDir()
}

// IsValidDirectory checks if the path is a directory.
func (v *pathValidator) IsValidDirectory(path string) bool {
	dirInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return dirInfo.IsDir()
}

// ConfigDirProvider is an interface that provides the config directory.
type ConfigDirProvider interface {
	// GetConfigDir gets the config directory.
	GetConfigDir() string
}

// ConfigDirResolver is a struct that resolves the config directory.
type ConfigDirResolver struct {
	PathValidator PathValidator
}

// NewConfigDirResolver creates a new ConfigDirResolver instance.
func NewConfigDirResolver(validator PathValidator) *ConfigDirResolver {
	return &ConfigDirResolver{PathValidator: validator}
}

// GetConfigDir gets the config directory.
func (c *ConfigDirResolver) GetConfigDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" && c.PathValidator.IsValidDirectory(os.Getenv("XDG_CONFIG_HOME")) {
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
	PathValidator     PathValidator
	ConfigDirProvider ConfigDirProvider
	Subdirs           []string
	Filename          string
}

// NewConfigFileResolver creates a new ConfigResolver instance.
func NewConfigFileResolver(validator PathValidator, provider ConfigDirProvider, subdirs []string, filename string) *ConfigFileResolver {
	return &ConfigFileResolver{
		PathValidator:     validator,
		ConfigDirProvider: provider,
		Subdirs:           subdirs,
		Filename:          filename,
	}
}

// Resolve resolves the config file path.
func (c *ConfigFileResolver) Resolve() (string, error) {
	for _, subdir := range c.Subdirs {
		fullPath := filepath.Join(c.ConfigDirProvider.GetConfigDir(), subdir)
		if !c.PathValidator.IsValidDirectory(fullPath) {
			slog.Debug("Config directory does not exist or is not a directory", "path", fullPath)
			continue
		}

		configPath := filepath.Join(fullPath, c.Filename)
		if c.PathValidator.IsValidFile(configPath) {
			return configPath, nil
		}
		slog.Debug("Config file not found", "path", configPath)
	}
	return "", fmt.Errorf("config file %s not found in any of the directories %v", c.Filename, c.Subdirs)
}
