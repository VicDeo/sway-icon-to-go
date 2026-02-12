package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// ConfigResolver is a struct that resolves the config file path.
type ConfigResolver struct {
	UserHomeDir string
	Path        string
	Filename    string
}

// NewConfigResolver creates a new ConfigResolver instance.
func NewConfigResolver(userHomeDir string, path string, filename string) *ConfigResolver {
	return &ConfigResolver{
		UserHomeDir: userHomeDir,
		Path:        path,
		Filename:    filename,
	}
}

// Resolve resolves the config file path.
func (c *ConfigResolver) Resolve() (string, error) {
	fullPath := filepath.Join(c.UserHomeDir, c.Path)
	if !isValidDirectory(fullPath) {
		return "", fmt.Errorf("directory %s does not exist or is not a directory", fullPath)
	}

	configPath := filepath.Join(fullPath, c.Filename)
	if !isValidFile(configPath) {
		return "", fmt.Errorf("file %s does not exist or is not a file", configPath)
	}
	return configPath, nil
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
