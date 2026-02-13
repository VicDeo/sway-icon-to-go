// Package config provides a way to get the config for the app
// including app name to icon mappings

package config

import (
	"log/slog"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

type IconToAppMap map[string][]string

type Config struct {
	AppIcons  IconToAppMap
	Length    int
	Delimiter string
	Uniq      bool
}

const (
	DefaultLength    = 12
	DefaultDelimiter = "|"
	DefaultUniq      = true
	NoMatch          = "_no_match"
	iconFileName     = "app-icons.yaml"
	faFileName       = "fa-icons.yaml"
)

var (
	configDirectories = []string{
		".config/sway",
		".config/i3", // legacy support for i3
	}

	defaultIcons = map[string]string{
		"chrome":      "\uf268",
		"comment":     "\uf075",
		"edit":        "\uf044",
		"folder-open": "\uf07c",
		"terminal":    "\uf120",
		"question":    "\uf128",
	}

	DefaultIconConfig = IconToAppMap{
		"firefox": []string{"firefox"},
		"chrome": []string{
			"chromium-browser",
			"chrome",
			"google-chrome",
		},
		"terminal": []string{
			"x-terminal-emulator",
			"XTerm",
			"konsole",
			"Konsole",
		},
		"shield-alt": []string{
			"keepassxc",
		},
		"cog": []string{
			"yast2",
		},
		"envelope": []string{
			"thunderbird",
		},
		"edit": []string{
			"jetbrains-idea-ce",
			"code",
			"cursor",
		},
		"folder-open": []string{
			"nautilus",
			".*krusader.*",
		},
		"music": []string{
			"clementine",
		},
		"play": []string{
			"vlc",
		},
		"comment": []string{
			"signal",
			"discord",
			"Telegram",
		},
		"question": []string{
			NoMatch,
		},
	}

	icons = defaultIcons
)

// NewConfig creates a new config for the app.
func NewConfig(delim string, uniq bool, length int, configPath string) (*Config, error) {
	iconConfig := DefaultIconConfig
	if configPath == "" {
		configPath = getConfigFilePath(iconFileName)
	}

	if configPath != "" {
		fileInfo, fileErr := os.Stat(configPath)
		if fileErr == nil && !fileInfo.IsDir() {
			viper.SetConfigFile(configPath)
			viper.SetConfigType("yaml")
			if err := viper.ReadInConfig(); err == nil {
				iconMap := &IconToAppMap{}
				err = viper.Unmarshal(iconMap)
				if err == nil {
					slog.Info("Config file is found", "path", configPath)
					iconConfig = *iconMap
				}
			}
		}
	}

	faIconsPath := getConfigFilePath(faFileName)
	if faIconsPath != "" {
		fileInfo, fileErr := os.Stat(faIconsPath)
		if fileErr == nil && !fileInfo.IsDir() {
			slog.Info("Font Awesome config file is found", "path", faIconsPath)
			viper.SetConfigFile(faIconsPath)
			viper.SetConfigType("yaml")
			if err := viper.ReadInConfig(); err == nil {
				faIcons := &map[string]string{}
				err = viper.Unmarshal(faIcons)
				if err == nil {
					slog.Info("Font Awesome config file is loaded", "path", faIconsPath)
					icons = *faIcons
					for k, v := range icons {
						icons[k], _ = strconv.Unquote(`"` + v + `"`)
					}
				}
			}
		}
	}

	currentConfig := &Config{
		AppIcons:  iconConfig,
		Length:    length,
		Delimiter: delim,
		Uniq:      uniq,
	}
	return currentConfig, nil
}

func (c *Config) GetAppIcon(name string) (string, bool) {
	// Note: we expect the name to be lowercase but this is the subject of a discussion
	name = strings.ToLower(name)

	for icon, appNames := range c.AppIcons {
		for _, appName := range appNames {
			match, err := regexp.MatchString(appName, name)
			if err != nil {
				slog.Warn("Error while matching app name", "name", name, "appName", appName, "error", err)
				continue
			}
			if match {
				return icons[icon], true
			}
		}
	}

	// TODO: make this configurable
	//return icons[NoMatch], false
	return name, false
}

func IsNoMatchIcon(icon string) bool {
	return icon == icons[NoMatch]
}

// getConfigFilePath gets the config file path for the given file name
func getConfigFilePath(fileName string) string {
	usr, _ := user.Current()
	home := usr.HomeDir
	for _, dir := range configDirectories {
		resolver := NewConfigResolver(home, dir, fileName)
		configPath, err := resolver.Resolve()
		if err != nil {
			slog.Warn("No config file found in", "directory", dir)
			continue
		}
		slog.Info("Config file found", "path", configPath)
		return configPath
	}
	return ""
}
