package main

import (
	"sway-icon-to-go/internal/config"
	"sway-icon-to-go/internal/proc"
	"sway-icon-to-go/internal/sway"
	"sync"
)

// IconProvider is a struct that provides the icon for the given pid and node name.
type IconProvider struct {
	processManager *proc.ProcessManager
	config         *config.Config
}

// NewIconProvider creates a new IconProvider instance.
func NewIconProvider(processManager *proc.ProcessManager, config *config.Config) *IconProvider {
	return &IconProvider{
		processManager: processManager,
		config:         config,
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

// GetIcon provides the icon for the given pid and node name.
func (i *IconProvider) GetIcon(pid *uint32, name string) (string, bool) {
	// Search by name first
	icon, found := config.GetAppIcon(name)
	if found {
		return icon, found
	}

	// Then search by pid
	appName, found := i.processManager.GetProcessName(pid)
	if !found {
		return name, found
	}

	icon, found = config.GetAppIcon(appName)
	if found {
		return icon, found
	}
	return name, found
}
