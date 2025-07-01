package sway

import (
	"fmt"
	"strings"
)

type Workspace struct {
	Name     string
	Number   int64
	AppIcons []string
}

func NewWorkspace(name string, number int64) *Workspace {
	return &Workspace{
		Name:     name,
		Number:   number,
		AppIcons: make([]string, 0),
	}
}

// TODO: use separator
func (w *Workspace) GetNewName() string {
	fmt.Println("Getting new name for workspace", w.AppIcons)
	return fmt.Sprintf("%d:%s", w.Number, strings.Join(w.AppIcons, ""))
}

func (w *Workspace) AddAppIcon(appIcon string) {
	w.AppIcons = append(w.AppIcons, appIcon)
}
