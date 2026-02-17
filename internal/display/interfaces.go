package display

// IconCache is an interface that provides a cache of icons.
type IconCache interface {
	GetIcon(name string) (string, bool)
	SetIcon(name string, icon string)
	Clear()
}

// ProcessManager is an interface that provides a process manager.
type ProcessManager interface {
	GetProcessName(pid *uint32) (string, bool)
}
