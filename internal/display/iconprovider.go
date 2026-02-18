// Package display provides a way to display the applications icons on the workspaces.

package display

import (
	"log/slog"
	"regexp"
	"strings"
	"sway-icon-to-go/internal/workspace"
	"sync"
)

// AppToIconMap is a map of application names to icons.
type AppToIconMap map[string]string

// IconProvider is a struct that provides the icon for the given pid and node name.
type IconProvider struct {
	processManager ProcessManager
	iconMap        AppToIconMap
	cache          IconCache
}

// NewIconProvider creates a new IconProvider instance.
func NewIconProvider(processManager ProcessManager, iconMap AppToIconMap, cache IconCache) *IconProvider {
	return &IconProvider{
		processManager: processManager,
		iconMap:        iconMap,
		cache:          cache,
	}
}

// AddIcons adds icons to the all windows of all workspaces.
func (i *IconProvider) AddIcons(workspaces workspace.Workspaces) error {
	var wg sync.WaitGroup
	for _, ws := range workspaces {
		wg.Add(1)
		go func(w *workspace.Workspace) {
			defer wg.Done()
			slog.Debug("Adding icons to workspace", "workspace", w.String())
			for _, window := range w.Windows {
				icon, found := i.GetIcon(window.PID, window.Title)
				if !found {
					icon = window.Title
				}
				w.AddAppIcon(icon)
			}
		}(ws)
	}
	wg.Wait()
	return nil
}

// ClearCache clears the cache.
func (i *IconProvider) ClearCache() {
	i.cache.Clear()
}

// GetIcon provides the icon for the given pid and node name.
func (i *IconProvider) GetIcon(pid *uint32, name string) (string, bool) {
	normalizedName := strings.ToLower(name)
	// Search by name first
	if icon, ok := i.iconFor(normalizedName); ok {
		return icon, true
	}

	// Then search by pid
	appName, ok := i.processManager.GetProcessName(pid)
	if !ok {
		return name, false
	}

	normalizedAppName := strings.ToLower(appName)
	if icon, ok := i.iconFor(normalizedAppName); ok {
		return icon, true
	}

	return name, false
}

func (i *IconProvider) iconFor(name string) (string, bool) {
	if icon, ok := i.cache.GetIcon(name); ok {
		return icon, true
	}

	icon, ok := i.iconMap[name]
	if ok {
		i.cache.SetIcon(name, icon)
		return icon, ok
	}

	// try treat app names in config as a regex and match them against the app name
	for appName, icon := range i.iconMap {
		if ok, err := regexp.MatchString(appName, name); err == nil && ok {
			i.cache.SetIcon(name, icon)
			return icon, true
		}
	}
	return name, false
}
