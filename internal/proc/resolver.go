package proc

import (
	"path/filepath"
	"strconv"
)

// LinuxResolver is a struct that resolves the process name by pid.
type LinuxResolver struct {
	ProcPath string
}

// Resolve resolves the process name by pid.
func (l *LinuxResolver) Resolve(pid uint32) (string, error) {
	exePath := filepath.Join(l.ProcPath, strconv.FormatUint(uint64(pid), 10), "exe")
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}
	return filepath.Base(realPath), nil
}
