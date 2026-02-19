// Package config provides a way to get the config for the app
// including app name to icon mappings.

package config

import (
	"log/slog"
	"os/user"
	"strconv"
	"strings"
)

// IconToAppMap is a map of icon names to one or more application names (regex patterns supported).
type IconToAppMap map[string][]string
type AppToIconMap map[string]string

// Config is a struct that contains the config for the app.
type Config struct {
	AppToIcon AppToIconMap
	Format    *Format
}

const (
	NoMatch      = "_no_match"
	iconFileName = "app-icons.yaml"
	faFileName   = "fa-icons.yaml"
)

var (
	configDirectories = []string{
		".config/sway",
		".config/i3", // legacy support for i3
	}

	defaultFaIcons = map[string]string{
		"chrome":      "\uf268",
		"comment":     "\uf075",
		"edit":        "\uf044",
		"folder-open": "\uf07c",
		"terminal":    "\uf120",
		"question":    "\uf128",
	}

	defaultIconConfig = IconToAppMap{
		"firefox": []string{"firefox"},
		"chrome": []string{
			"chromium-browser",
			"chrome",
			"google-chrome",
		},
		"terminal": []string{
			"x-terminal-emulator",
			"xterm",
			"konsole",
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
		"question-circle": []string{
			NoMatch,
		},
	}
)

// NewConfig creates a new config for the app.
func NewConfig(configPath string, format *Format) (*Config, error) {
	if format == nil {
		format = DefaultFormat()
		slog.Warn("No format provided, using default format")
	}
	iconConfig := defaultIconConfig
	faIcons := defaultFaIcons

	if configPath == "" {
		configPath = getConfigFilePath(iconFileName)
	}

	if configPath != "" {
		configFile, err := NewConfigFile(configPath)
		// if error just use default icons
		if err == nil {
			loadedIconConfig := &IconToAppMap{}
			if err := configFile.Load(loadedIconConfig); err == nil {
				iconConfig = *loadedIconConfig
			}
		}
	}

	faIconsPath := getConfigFilePath(faFileName)
	if faIconsPath != "" {
		configFile, err := NewConfigFile(faIconsPath)
		// if error just use default Font Awesome icons
		if err == nil {
			loadedFaIcons := &map[string]string{}
			if err := configFile.Load(loadedFaIcons); err == nil {
				faIcons = *loadedFaIcons
				for k, v := range faIcons {
					var err error
					// Hack: quote and unquote to get the unicode character from code
					faIcons[k], err = strconv.Unquote(`"` + v + `"`)
					if err != nil {
						slog.Error("Error while unquoting icon", "icon", v, "error", err)
					}
				}
			}
		}
	}

	// Now transform icon name []app name to
	// icon name [app name1, app name2, ...] to
	// app name1[icon name], app name2[icon name], ...
	iconByAppName := make(map[string]string)
	for icon, appNames := range iconConfig {
		faIcon, ok := faIcons[icon]
		if !ok {
			slog.Warn("FA icon not found", "icon", icon)
			continue
		}

		for _, appName := range appNames {
			// Note: we expect the name to be lowercase but this is the subject of a discussion
			iconByAppName[strings.ToLower(appName)] = faIcon
		}
	}

	currentConfig := &Config{
		AppToIcon: iconByAppName,
		Format:    format,
	}
	return currentConfig, nil
}

// getConfigFilePath gets the config file path for the given file name
func getConfigFilePath(fileName string) string {
	usr, err := user.Current()
	if err != nil {
		slog.Error("Error while getting current user", "error", err)
		return ""
	}
	home := usr.HomeDir
	for _, dir := range configDirectories {
		resolver := NewConfigResolver(home, dir, fileName)
		configPath, err := resolver.Resolve()
		if err != nil {
			slog.Warn("No config file found", "directory", dir, "error", err)
			continue
		}
		slog.Info("Config file found", "path", configPath)
		return configPath
	}
	return ""
}
