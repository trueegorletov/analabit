package cache

import (
	"encoding/json"
	"os"
	"sync"
)

// LayeredCache provides two-level lookup: an in-memory TTL cache (fast) backed by a persistent
// store (slow). The design keeps the interface intentionally simple while the persistent layer is
// evolving – it can be swapped for a real database implementation later without touching callers.

// PersistentStore is the minimal contract required from a slower persistent layer. Implementations
// should be concurrency-safe.
type PersistentStore interface {
	Get(key string) (any, bool)
	Set(key string, val any)
}

// LayeredCache first consults the in-memory cache. On a miss it falls back to the persistent
// store; if found there, the value is promoted back to the memory layer for future fast access.
// All writes go through to both layers.
//
// IMPORTANT: TTL is enforced only in the memory layer – persistent store keeps the value
// indefinitely (or according to its own retention policies).
type LayeredCache struct {
	mem        *MemoryCache
	persistent PersistentStore
}

// NewLayeredCache wires an existing MemoryCache with a persistent store.
func NewLayeredCache(mem *MemoryCache, ps PersistentStore) *LayeredCache {
	return &LayeredCache{mem: mem, persistent: ps}
}

// Get tries memory first, then persistent. If the value is found in persistent the method will
// promote it back to the memory layer.
func (lc *LayeredCache) Get(key string) (any, bool) {
	if v, ok := lc.mem.Get(key); ok {
		return v, true
	}
	if v, ok := lc.persistent.Get(key); ok {
		// promote to fast layer for future hits
		lc.mem.Set(key, v)
		return v, true
	}
	return nil, false
}

// Set writes to both layers. Persistent layer is considered the source of truth – if it fails
// the error is ignored for now (future versions may surface it).
func (lc *LayeredCache) Set(key string, val any) {
	lc.mem.Set(key, val)
	lc.persistent.Set(key, val)
}

// Purge delegates to the memory layer. Persistent store purging (if any) is left to its own
// maintenance routines.
func (lc *LayeredCache) Purge() { lc.mem.Purge() }

// --- Simple map-based persistent store for early tests ------------------------------------------------------------

// MapStore is a naïve in-process implementation of PersistentStore backed by a mutex-protected map.
// It has no TTL and is NOT intended for production.
type MapStore struct {
	mu       sync.RWMutex
	data     map[string]any
	filePath string
}

func NewMapStore(filePath string) *MapStore {
	m := &MapStore{data: make(map[string]any), filePath: filePath}
	m.Load()
	return m
}

func (m *MapStore) Get(key string) (any, bool) {
	m.mu.RLock()
	v, ok := m.data[key]
	m.mu.RUnlock()
	return v, ok
}

func (m *MapStore) Set(key string, val any) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = val

	// Inlined Save() logic to persist on every set
	file, err := os.Create(m.filePath)
	if err != nil {
		// As per original design, persistence errors are not surfaced.
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	_ = encoder.Encode(m.data) // Error is ignored
}

func (m *MapStore) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	file, err := os.Open(m.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Not an error if file doesn't exist yet
		}
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&m.data)
}
