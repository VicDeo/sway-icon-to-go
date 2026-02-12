// Package sway provides a way to process the workspaces and renames them
// according to the name formatter and icon provider basing on the apps running on the workspaces.

package sway

import (
	"context"
)

// IconProvider is an interface that provides an icon for a given PID and node name.
type IconProvider interface {
	AddIcons(workspaces Workspaces) error
}

// NameFormatter is an interface that formats a workspace name.
type NameFormatter interface {
	Format(workspaceNumber int64, appIcons []string) string
}

// ProcessWorkspaces processes the workspaces and renames them according to the name formatter and icon provider
// basing on the apps running on the workspaces.
func ProcessWorkspaces(ctx context.Context, iconProvider IconProvider, nameFormatter NameFormatter) error {
	sway, err := NewSwayClient(ctx)
	if err != nil {
		return err
	}

	// Then traverse the tree and populate the workspaces map.
	workspaces, err := sway.CollectWorkspaces()
	if err != nil {
		return err
	}

	// Add icons to the all windows of all workspaces.
	if err := iconProvider.AddIcons(workspaces); err != nil {
		return err
	}

	// Send all commands at once as there could be a mess otherwise.
	if err := sway.RenameWorkspaces(workspaces, nameFormatter); err != nil {
		return err
	}
	return nil
}
