package proc

import (
	"log/slog"
	"sync"
	"time"
)

const (
	// Set the cache TTL to 15 minutes.
	CacheTTL = 15 * time.Minute
)

// NameResolver resolves the name of the process executable
type NameResolver interface {
	Resolve(pid uint32) (string, error)
}

// ProcessManager holds cache of filename by pid
type ProcessManager struct {
	mu       sync.Mutex
	cache    map[uint32]CacheItem
	resolver NameResolver
}

// NewProcessManager creates ProcessManager instance
func NewProcessManager(resolver NameResolver) *ProcessManager {
	return &ProcessManager{
		cache:    make(map[uint32]CacheItem),
		resolver: resolver,
	}
}

type CacheItem struct {
	Name       string
	BestBefore time.Time
}

// GetProcessName gets process name by pid
func (pm *ProcessManager) GetProcessName(pid *uint32) (string, bool) {
	if pid == nil {
		return "", false
	}
	pm.mu.Lock()
	if item, ok := pm.cache[*pid]; ok && item.BestBefore.After(time.Now()) {
		pm.mu.Unlock()
		return item.Name, true
	}
	pm.mu.Unlock()

	name, err := pm.resolver.Resolve(*pid)
	if err != nil {
		slog.Error("error while getting executable name", "error", err)
		return "", false
	}

	pm.mu.Lock()
	pm.cache[*pid] = CacheItem{Name: name, BestBefore: time.Now().Add(CacheTTL)}
	pm.mu.Unlock()

	return name, true
}
