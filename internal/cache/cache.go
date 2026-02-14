package cache

import (
	"sync"
)

const (
	// 30 different applications is a reasonable default for the cache capacity.
	InitialCacheCapacity = 30
)

// Cache is a struct that caches the icons.
type Cache struct {
	nameCache map[string]string
	muName    sync.Mutex
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return &Cache{
		nameCache: make(map[string]string, InitialCacheCapacity),
		muName:    sync.Mutex{},
	}
}

// Clear clears the cache.
func (c *Cache) Clear() {
	c.muName.Lock()
	defer c.muName.Unlock()
	c.nameCache = make(map[string]string, InitialCacheCapacity)
}

// GetIcon gets the icon for the given name.
func (c *Cache) GetIcon(name string) (string, bool) {
	c.muName.Lock()
	defer c.muName.Unlock()
	icon, ok := c.nameCache[name]
	return icon, ok
}

// SetIcon sets the icon for the given name.
func (c *Cache) SetIcon(name string, icon string) {
	c.muName.Lock()
	defer c.muName.Unlock()
	c.nameCache[name] = icon
}
