package sway

// Workspace is a struct that represents a workspace
type Workspace struct {
	Name     string
	Number   int64
	AppIcons []string
}

// NewWorkspace creates a new workspace
func NewWorkspace(name string, number int64) *Workspace {
	return &Workspace{
		Name:     name,
		Number:   number,
		AppIcons: make([]string, 0, 10),
	}
}

// GetNewName gets the new name for the workspace
func (w *Workspace) GetNewName(nameFormatter NameFormatter) string {
	return nameFormatter.Format(w.Number, w.AppIcons)
}

// AddAppIcon adds an app icon to the workspace
func (w *Workspace) AddAppIcon(appIcon string) {
	w.AppIcons = append(w.AppIcons, appIcon)
}
