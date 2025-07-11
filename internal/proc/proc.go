package proc

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	procPath         = "/proc"
	DefaultCacheSize = 30
)

var (
	pidCache = make(map[uint32]string, DefaultCacheSize)
	mu       sync.Mutex
)

func GetProcessName(pid *uint32) (string, bool) {
	mu.Lock()
	defer mu.Unlock()
	if pid == nil {
		return "", false
	}
	if cachedName, ok := pidCache[*pid]; ok {
		return cachedName, true
	}
	filename, err := getExecutableName(pid)
	if err != nil || filename == "" {
		log.Printf("Error while getting executable name: %v\n", err)
		return "", false
	}
	pidCache[*pid] = filename
	return filename, true
}

// Resolve the executable name for the given pid
func getExecutableName(pid *uint32) (string, error) {
	if pid == nil {
		return "", fmt.Errorf("pid is nil")
	}
	pidStr := strconv.FormatUint(uint64(*pid), 10)
	exePath := filepath.Join(procPath, pidStr, "exe")
	realPath, err := filepath.EvalSymlinks(exePath)
	if err != nil {
		return "", err
	}
	return filepath.Base(realPath), nil
}
