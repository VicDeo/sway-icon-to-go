package config

import (
	"fmt"
	"regexp"
	"strings"
)

func GetAppIcon(name string) (string, bool) {
	// Note:we expect the name to be lowercase but this is the subject of a discussion
	name = strings.ToLower(name)
	if iconCache[name] != "" {
		return iconCache[name], true
	}

	for icon, appNames := range currentConfig.AppIcons {
		for _, appName := range appNames {
			match, _ := regexp.MatchString(appName, name)
			if match {
				iconCache[name] = icons[icon]
				return iconCache[name], true
			}
		}
	}

	// TODO: make this configurable
	iconCache[name] = TrimAppName(name)
	//iconCache[name] = icons[NoMatch]

	return iconCache[name], false
}

func IsNoMatchIcon(icon string) bool {
	return icon == icons[NoMatch]
}

func TrimAppName(appName string) string {
	if len(appName) > currentConfig.Length {
		return appName[:currentConfig.Length]
	}
	return appName
}

func BuildName(id int64, appNames []string) string {
	if currentConfig.Uniq {
		appNames = unique(appNames)
	}
	return fmt.Sprintf(
		"%d: %s",
		id,
		strings.Join(appNames, currentConfig.Delimiter),
	)
}

func unique(slice []string) []string {
	var unique []string

	uniqueMap := map[string]int{}
	for _, v := range slice {
		if _, exist := uniqueMap[v]; !exist {
			uniqueMap[v] = 1
			unique = append(unique, v)
		} else {
			uniqueMap[v]++
		}
	}
	return unique
}
