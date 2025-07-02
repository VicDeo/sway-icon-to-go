package sway

import (
	"fmt"
	"strings"
)

// getRenameWorkspaceCommand gets the command to rename the workspace
func getRenameWorkspaceCommand(oldName string, newName string) string {
	return fmt.Sprintf(
		"rename workspace \"%s\" to \"%s\"",
		escapeName(oldName),
		escapeName(newName),
	)
}

// getMoveAppToWorkspaceCommand gets the command to move the app to the workspace
func getMoveAppToWorkspaceCommand(appId int64, workspaceName string) string {
	return fmt.Sprintf(
		"[con_id=\"%d\"] move container to workspace \"%s\"",
		appId,
		escapeName(workspaceName),
	)
}

// getFocusCommand gets the command to focus the workspace
func getFocusCommand(workspaceName string) string {
	return fmt.Sprintf("workspace \"%s\"", escapeName(workspaceName))
}

// escapeName escapes the name for the sway command
func escapeName(name string) string {
	return strings.Replace(name, "\"", "\\\"", -1)
}
