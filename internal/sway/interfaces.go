package sway

// NameFormatter is an interface that formats a workspace name.
type NameFormatter interface {
	Format(workspaceNumber int64, appIcons []string) string
}
