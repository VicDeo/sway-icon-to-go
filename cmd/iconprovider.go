package main

import (
	"regexp"
	"sway-icon-to-go/internal/cache"
	"sway-icon-to-go/internal/config"
	"sway-icon-to-go/internal/proc"
	"sway-icon-to-go/internal/sway"
	"sync"
)

// IconProvider is a struct that provides the icon for the given pid and node name.
type IconProvider struct {
	processManager *proc.ProcessManager
	config         *config.Config
	cache          *cache.Cache
}

// NewIconProvider creates a new IconProvider instance.
func NewIconProvider(processManager *proc.ProcessManager, config *config.Config, cache *cache.Cache) *IconProvider {
	return &IconProvider{
		processManager: processManager,
		config:         config,
		cache:          cache,
	}
}

// AddIcons adds icons to the all windows of all workspaces.
func (i *IconProvider) AddIcons(workspaces sway.Workspaces) error {
	var wg sync.WaitGroup
	for wn, workspace := range workspaces {
		wg.Add(1)
		go func(wn int64, workspace *sway.Workspace) {
			defer wg.Done()
			for _, window := range workspace.Windows {
				icon, found := i.GetIcon(window.PID, window.Title)
				if !found {
					icon = window.Title
				}
				workspace.AddAppIcon(icon)
			}
		}(wn, workspace)
	}
	wg.Wait()
	return nil
}

func (i *IconProvider) ClearCache() {
	i.cache.Clear()
}

// GetIcon provides the icon for the given pid and node name.
func (i *IconProvider) GetIcon(pid *uint32, name string) (string, bool) {

	if icon, ok := i.cache.GetIcon(name); ok {
		return icon, true
	}

	// Search by name first
	if icon, ok := i.iconFor(name); ok {
		i.cache.SetIcon(name, icon)
		return icon, true
	}

	// Then search by pid
	appName, ok := i.processManager.GetProcessName(pid)
	if !ok {
		return name, false
	}

	if icon, ok := i.iconFor(appName); ok {
		i.cache.SetIcon(appName, icon)
		return icon, true
	}

	return name, false
}

func (i *IconProvider) iconFor(name string) (string, bool) {
	icon, found := i.config.AppToIcon[name]
	if found {
		return icon, found
	}

	// try treat app names in config as a regex and match them against the app name
	for appName, icon := range i.config.AppToIcon {
		if ok, err := regexp.MatchString(appName, name); err == nil && ok {
			return icon, true
		}
	}
	return name, false
}
