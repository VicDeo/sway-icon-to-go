package sway

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

func (w *Workspace) GetNewName(nameFormatter NameFormatter) string {
	return nameFormatter.Format(w.Number, w.AppIcons)
}

func (w *Workspace) AddAppIcon(appIcon string) {
	w.AppIcons = append(w.AppIcons, appIcon)
}
