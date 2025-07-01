package sway

import (
	"fmt"
	"strings"
)

func getRenameWorkspaceCommand(oldName string, newName string) string {
	return fmt.Sprintf(
		"rename workspace \"%s\" to \"%s\"",
		escapeName(oldName),
		escapeName(newName),
	)
}

func getMoveAppToWorkspaceCommand(appId int64, workspaceName string) string {
	return fmt.Sprintf(
		"[con_id=\"%d\"] move container to workspace \"%s\"",
		appId,
		escapeName(workspaceName),
	)
}

func getFocusCommand(workspaceName string) string {
	return fmt.Sprintf("workspace \"%s\"", escapeName(workspaceName))
}

func escapeName(name string) string {
	return strings.Replace(name, "\"", "\\\"", -1)
}
