package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

type ConfigLoader struct {
	viper *viper.Viper
	Path  string
}

func NewConfigLoader(path string) (*ConfigLoader, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		slog.Error("Error while statting config file", "error", err)
		return nil, fmt.Errorf("error while getting file info for %s: %w", path, err)
	}
	if fileInfo.IsDir() {
		slog.Error("Config file is a directory", "path", path)
		return nil, fmt.Errorf("config file %s is not a file", path)
	}
	return &ConfigLoader{
		viper: viper.New(),
		Path:  path,
	}, nil
}

func (c *ConfigLoader) Load(target any) error {
	c.viper.SetConfigFile(c.Path)
	c.viper.SetConfigType("yaml")
	if err := c.viper.ReadInConfig(); err != nil {
		slog.Error("Error while reading config file", "error", err)
		return fmt.Errorf("error while reading config file %s: %w", c.Path, err)
	}

	if err := c.viper.Unmarshal(target); err != nil {
		slog.Error("Error while unmarshalling config file", "error", err)
		return fmt.Errorf("error while unmarshalling config file %s: %w", c.Path, err)
	}

	slog.Debug("Config file is loaded successfully", "path", c.Path)

	return nil
}
