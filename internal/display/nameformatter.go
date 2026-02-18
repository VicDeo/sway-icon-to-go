package display

import (
	"fmt"
	"strings"
)

// Format is a struct that contains the format params for the workspace name.
type Format struct {
	Length    int
	Delimiter string
	Uniq      bool
}

// NameFormatter is a struct that formats the workspace name according to the config.
type NameFormatter struct {
	format *Format
}

// NewNameFormatter creates a new NameFormatter with the given config.
func NewNameFormatter(delimiter string, length int, uniq bool) *NameFormatter {
	format := &Format{
		Delimiter: delimiter,
		Length:    length,
		Uniq:      uniq,
	}
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

	trimmedAppIcons := make([]string, 0)
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
