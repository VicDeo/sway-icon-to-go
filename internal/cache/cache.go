// Package cache provides a cache for the icons.

package cache

import (
	"slices"
	"sync"
)

const (
	// 50 different applications is a reasonable default for the cache capacity.
	InitialCacheCapacity = 50
)

// Cache is a struct that caches the icons.
type Cache struct {
	nameCache   map[string]string
	order       []string
	muName      sync.Mutex
	maxCapacity int
}

// NewCache creates a new Cache instance.
func NewCache() *Cache {
	return NewCappedCache(InitialCacheCapacity)
}

// NewCappedCache creates a new Cache instance with a maximum capacity.
func NewCappedCache(maxCapacity int) *Cache {
	if maxCapacity <= 0 {
		maxCapacity = InitialCacheCapacity
	}
	return &Cache{
		nameCache:   make(map[string]string, maxCapacity),
		order:       make([]string, 0, maxCapacity),
		maxCapacity: maxCapacity,
	}
}

// Clear clears the cache.
func (c *Cache) Clear() {
	c.muName.Lock()
	defer c.muName.Unlock()
	c.nameCache = make(map[string]string, c.maxCapacity)
	c.order = make([]string, 0, c.maxCapacity)
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
	if _, exists := c.nameCache[name]; exists {
		// move to the end of the order
		c.order = slices.Delete(c.order, slices.Index(c.order, name), 1)
		c.order = append(c.order, name)
		c.nameCache[name] = icon
		return
	}
	// evict oldest
	if len(c.nameCache) >= c.maxCapacity {
		oldest := c.order[0]
		delete(c.nameCache, oldest)
		c.order = c.order[1:]
	}
	c.order = append(c.order, name)
	c.nameCache[name] = icon
}
