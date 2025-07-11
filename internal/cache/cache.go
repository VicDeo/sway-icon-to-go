package cache

import (
	"sync"
)

const (
	DefaultCacheSize = 30
)

var (
	nameCache = make(map[string]string, DefaultCacheSize)
	muName    sync.Mutex
)

func Clear() {
	muName.Lock()
	defer muName.Unlock()
	nameCache = make(map[string]string, DefaultCacheSize)
}

func GetIcon(name string) (string, bool) {
	muName.Lock()
	defer muName.Unlock()
	icon, ok := nameCache[name]
	return icon, ok
}

func SetIcon(name string, icon string) {
	muName.Lock()
	defer muName.Unlock()
	nameCache[name] = icon
}
