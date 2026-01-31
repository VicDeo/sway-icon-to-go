package proc

import (
	"log"
	"sync"
)

// NameResolver resolves the name of the process executable
type NameResolver interface {
	Resolve(pid uint32) (string, error)
}

// ProcessManager holds cache of filename by pid
type ProcessManager struct {
	mu       sync.Mutex
	cache    map[uint32]string
	resolver NameResolver
}

// NewProcessManager creates ProcessManager instance
func NewProcessManager(resolver NameResolver) *ProcessManager {
	return &ProcessManager{
		cache:    make(map[uint32]string),
		resolver: resolver,
	}
}

// GetProcessName gets process name by pid
func (pm *ProcessManager) GetProcessName(pid *uint32) (string, bool) {
	if pid == nil {
		return "", false
	}
	pm.mu.Lock()
	if name, ok := pm.cache[*pid]; ok {
		pm.mu.Unlock()
		return name, true
	}
	pm.mu.Unlock()

	name, err := pm.resolver.Resolve(*pid)
	if err != nil {
		log.Printf("error while getting executable name: %v\n", err)
		return "", false
	}

	pm.mu.Lock()
	pm.cache[*pid] = name
	pm.mu.Unlock()

	return name, true
}
