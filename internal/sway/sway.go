// Package sway provides a way to process the workspaces and renames them
// according to the name formatter and icon provider basing on the apps running on the workspaces

package sway

import (
	"context"
	"log"
)

// IconProvider is an interface that provides an icon for a given PID and node name
type IconProvider interface {
	GetIcon(pid *uint32, nodeName string) (string, bool)
}

// NameFormatter is an interface that formats a workspace name
type NameFormatter interface {
	Format(workspaceNumber int64, appIcons []string) string
}

// ProcessWorkspaces processes the workspaces and renames them according to the name formatter and icon provider
// basing on the apps running on the workspaces
func ProcessWorkspaces(ctx context.Context, iconProvider IconProvider, nameFormatter NameFormatter) error {
	sway, err := NewSwayClient(ctx)
	if err != nil {
		return err
	}

	// Then traverse the tree and populate the workspaces map
	workspaces, err := sway.CollectWorkspaces()
	if err != nil {
		return err
	}

	// Populate the workspaces with the app icons
	for _, workspace := range workspaces {
		for _, window := range workspace.Windows {
			icon, found := iconProvider.GetIcon(window.PID, window.Title)
			if !found {
				log.Printf("No app mapping found for %+v\n", window)
			}
			workspaces[workspace.Number].AddAppIcon(icon)
		}
	}

	// Send all commands at once as there could be a mess otherwise
	if error := sway.RenameWorkspaces(workspaces, nameFormatter); error != nil {
		return error
	}
	return nil
}
