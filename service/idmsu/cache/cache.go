package cache

import (
	"sync"
	"time"
)

// MemoryCache is a lightweight TTL cache used while the persistent layer is under development.
// Keys are arbitrary strings (e.g., program+competition identifiers) and values are opaque (interface{}).
// It is safe for concurrent use.

type MemoryCache struct {
	mu   sync.RWMutex
	ttl  time.Duration
	data map[string]entry
}

type entry struct {
	val    any
	expiry time.Time
}

// NewMemoryCache constructs a cache with the given TTL.
func NewMemoryCache(ttl time.Duration) *MemoryCache {
	return &MemoryCache{
		ttl:  ttl,
		data: make(map[string]entry),
	}
}

// Get retrieves a value if present and not expired.
func (c *MemoryCache) Get(key string) (any, bool) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expiry) {
		return nil, false
	}
	return e.val, true
}

// Set stores a value with TTL.
func (c *MemoryCache) Set(key string, val any) {
	c.mu.Lock()
	c.data[key] = entry{val: val, expiry: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}

// Purge removes expired items; can be called periodically.
func (c *MemoryCache) Purge() {
	now := time.Now()
	c.mu.Lock()
	for k, e := range c.data {
		if now.After(e.expiry) {
			delete(c.data, k)
		}
	}
	c.mu.Unlock()
}
