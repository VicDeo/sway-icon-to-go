package sway

import (
	"fmt"
	"log/slog"
	"strings"
)

// Workspace is a struct that represents a workspace
type Workspace struct {
	Name     string
	Number   int64
	Windows  []WindowInfo
	AppIcons []string
}

type WindowInfo struct {
	PID   *uint32
	Title string
}

// NewWorkspace creates a new workspace
func NewWorkspace(name string, number int64) *Workspace {
	return &Workspace{
		Name:     name,
		Number:   number,
		Windows:  make([]WindowInfo, 0, 10),
		AppIcons: make([]string, 0, 10),
	}
}

// AddAppIcon adds an app icon to the workspace
func (w *Workspace) AddAppIcon(appIcon string) {
	w.AppIcons = append(w.AppIcons, appIcon)
}

// AddWindow adds a window to the workspace
func (w *Workspace) AddWindow(window WindowInfo) {
	w.Windows = append(w.Windows, window)
}

// ToRenameCommand produces Sway rename command for the workspace
func (w *Workspace) ToRenameCommand(nf NameFormatter) string {
	newName := nf.Format(w.Number, w.AppIcons)
	// Do not rename if nothing has been changed
	if newName == w.Name {
		return ""
	}

	return fmt.Sprintf(
		"rename workspace \"%s\" to \"%s\"",
		escapeName(w.Name),
		escapeName(newName),
	)
}

// escapeName escapes the name for the sway command
func escapeName(name string) string {
	return strings.ReplaceAll(name, "\"", "\\\"")
}

// Workspaces is a map of workspace number to workspace
type Workspaces map[int64]*Workspace

// ToRenameCommand produces Sway rename command for all workspaces
func (ww *Workspaces) ToRenameCommand(nf NameFormatter) string {
	var commands []string

	for _, workspace := range *ww {
		name := workspace.ToRenameCommand(nf)
		if name != "" {
			commands = append(commands, name)
		}
	}
	command := strings.Join(commands, ";")
	slog.Info("Command is ready to be executed", "command", command)
	return command
}
