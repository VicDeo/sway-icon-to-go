package display

import (
	"fmt"
	"strings"
	"sway-icon-to-go/internal/config"
)

// NameFormatter is a struct that formats the workspace name according to the config.
type NameFormatter struct {
	format *config.Format
}

// NewNameFormatter creates a new NameFormatter with the given config.
func NewNameFormatter(format *config.Format) *NameFormatter {
	return &NameFormatter{format: format}
}

// Format the workspace name according to the config.
func (nf *NameFormatter) Format(workspaceNumber int64, appIcons []string) string {
	if appIcons == nil {
		return ""
	}

	if nf.format.Uniq {
		appIcons = unique(appIcons)
	}

	trimmedAppIcons := []string{}
	if nf.format.Length > 0 {
		// Trim app icons to the length specified in the config.
		for _, appIcon := range appIcons {
			capLength := min(len(appIcon), nf.format.Length)
			trimmedAppIcons = append(trimmedAppIcons, appIcon[:capLength])
		}
	} else {
		trimmedAppIcons = appIcons
	}
	return fmt.Sprintf(
		"%d: %s",
		workspaceNumber,
		strings.Join(trimmedAppIcons, nf.format.Delimiter),
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
