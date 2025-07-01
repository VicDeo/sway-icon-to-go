package sway

import (
	"strings"

	swayClient "github.com/joshuarubin/go-sway"
)

type UnifiedWorkspace struct {
	Workspace swayClient.Workspace
	Node      *swayClient.Node
	MergeWith *swayClient.Workspace
}

// Helper methods for UnifiedWorkspace
func (w UnifiedWorkspace) GetNewName() string {
	// This would need to be implemented based on your naming logic
	// For now, return a placeholder
	return w.Workspace.Name
}

func (w UnifiedWorkspace) IsCustom() bool {
	// Check if workspace name contains a colon (custom workspace)
	return strings.Contains(w.Workspace.Name, ":")
}

func (w UnifiedWorkspace) IsDuplicateOf(ww UnifiedWorkspace) bool {
	return w.getIdByName() == ww.getIdByName()
}

func (w UnifiedWorkspace) getIdByName() string {
	if idx := strings.IndexByte(w.Workspace.Name, ':'); idx >= 0 {
		return w.Workspace.Name[:idx]
	}
	return w.Workspace.Name
}
