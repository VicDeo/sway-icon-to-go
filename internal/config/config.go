// Package config provides a way to get the config for the app
// including app name to icon mappings

package config

import (
	"log"
	"os"
	"os/user"
	"strconv"

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

	icons         = defaultIcons
	currentConfig = &Config{}
)

// GetConfig creates a new config for the app
func GetConfig(delim string, uniq bool, length int, configPath string) (*Config, error) {
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
					log.Printf("Config file is found  %s\n", configPath)
					iconConfig = *iconMap
				}
			}
		}
	}

	faIconsPath := getConfigFilePath(faFileName)
	if faIconsPath != "" {
		fileInfo, fileErr := os.Stat(faIconsPath)
		if fileErr == nil && !fileInfo.IsDir() {
			log.Printf("Font Awesome config file is found  %s\n", faIconsPath)
			viper.SetConfigFile(faIconsPath)
			viper.SetConfigType("yaml")
			if err := viper.ReadInConfig(); err == nil {
				faIcons := &map[string]string{}
				err = viper.Unmarshal(faIcons)
				if err == nil {
					log.Printf("Font Awesome config file is loaded %s\n", faIconsPath)
					icons = *faIcons
					for k, v := range icons {
						icons[k], _ = strconv.Unquote(`"` + v + `"`)
					}
				}
			}
		}
	}
	currentConfig.AppIcons = iconConfig
	currentConfig.Length = length
	currentConfig.Delimiter = delim
	currentConfig.Uniq = uniq
	return currentConfig, nil
}

// getConfigFilePath gets the config file path for the given file name
func getConfigFilePath(fileName string) string {
	usr, _ := user.Current()
	home := usr.HomeDir
	for _, dir := range configDirectories {
		resolver := NewConfigResolver(home, dir, fileName)
		configPath, err := resolver.Resolve()
		if err != nil {
			log.Println("No config file found in", dir)
			continue
		}
		log.Println("Config file found in", configPath)
		return configPath
	}
	return ""
}
