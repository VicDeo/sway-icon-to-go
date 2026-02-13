package main

import (
	"fmt"
	"strings"
	"sway-icon-to-go/internal/config"
)

// NameFormatter is a struct that formats the workspace name according to the config.
type NameFormatter struct {
	config *config.Config
}

// NewNameFormatter creates a new ConfigNameFormatter with the given config.
func NewNameFormatter(config *config.Config) *NameFormatter {
	return &NameFormatter{config: config}
}

// Format the workspace name according to the config.
func (c NameFormatter) Format(workspaceNumber int64, appIcons []string) string {
	if appIcons == nil {
		return ""
	}

	if c.config.Uniq {
		appIcons = unique(appIcons)
	}

	trimmedAppIcons := make([]string, 0)
	if c.config.Length > 0 {
		// Trim app icons to the length specified in the config.
		for _, appIcon := range appIcons {
			capLength := min(len(appIcon), c.config.Length)
			trimmedAppIcons = append(trimmedAppIcons, appIcon[:capLength])
		}
	} else {
		trimmedAppIcons = appIcons
	}
	return fmt.Sprintf(
		"%d: %s",
		workspaceNumber,
		strings.Join(trimmedAppIcons, c.config.Delimiter),
	)
}

func unique(slice []string) []string {
	var uniqueApps []string

	seen := make(map[string]struct{})
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			uniqueApps = append(uniqueApps, v)
		}
	}
	return uniqueApps
}
