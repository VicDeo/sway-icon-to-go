package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
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
		".i3",
		".config/i3",
		".config/i3-regolith",
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
	iconCache     = map[string]string{}
	currentConfig = &Config{}
)

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
					//fmt.Println(configPath + " is loaded")
					iconConfig = *iconMap
				}
			}
		}
	}

	faIconsPath := getConfigFilePath(faFileName)
	if faIconsPath != "" {
		fileInfo, fileErr := os.Stat(faIconsPath)
		if fileErr == nil && !fileInfo.IsDir() {
			fmt.Println(faIconsPath)
			viper.SetConfigFile(faIconsPath)
			viper.SetConfigType("yaml")
			if err := viper.ReadInConfig(); err == nil {
				faIcons := &map[string]string{}
				err = viper.Unmarshal(faIcons)
				if err == nil {
					fmt.Println(faIconsPath + " is loaded")
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

func getConfigFilePath(fileName string) string {
	usr, _ := user.Current()
	home := usr.HomeDir
	for _, dir := range configDirectories {
		fullDir := filepath.Join(home, dir)
		dirInfo, dirErr := os.Stat(fullDir)
		if dirErr != nil || !dirInfo.IsDir() {
			continue
		}
		configPath := filepath.Join(fullDir, fileName)
		fileInfo, fileErr := os.Stat(configPath)
		if fileErr != nil || fileInfo.IsDir() {
			continue
		}
		return configPath
	}
	return ""
}
