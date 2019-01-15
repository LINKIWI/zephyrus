package cache

import (
	"sync"
	"time"
)

// MemoryTTLCache is an in-memory key-value cache with TTL support.
type MemoryTTLCache struct {
	// Underlying data store.
	store map[string]*cacheEntry
	// Mutex used to synchronize reads and writes to the store.
	mutex sync.Mutex
}

// cacheEntry is an internal data structure to represent an item in the cache.
type cacheEntry struct {
	value  interface{}
	expiry time.Time
}

// NewMemoryTTLCache creates a new MemoryTTLCache with default options.
func NewMemoryTTLCache() *MemoryTTLCache {
	return &MemoryTTLCache{
		store: make(map[string]*cacheEntry),
	}
}

// Get retrieves the value associated with a key. Returns the value if present and nil otherwise.
func (m *MemoryTTLCache) Get(key string) interface{} {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	entry, ok := m.store[key]

	if !ok {
		return nil
	}

	if entry.isExpired() {
		delete(m.store, key)
		return nil
	}

	return entry.value
}

// Set writes a new or updated key-value pair to the cache, with the specified TTL.
// Specify 0 as the TTL to never expire the cache entry.
func (m *MemoryTTLCache) Set(key string, value interface{}, ttl time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	expiry := time.Unix(0, 0)
	if ttl != 0 {
		expiry = time.Now().Add(ttl)
	}

	m.store[key] = &cacheEntry{
		value:  value,
		expiry: expiry,
	}
}

// Delete deletes an entry from the cache. Returns true if an item was deleted; false otherwise.
func (m *MemoryTTLCache) Delete(key string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	_, ok := m.store[key]

	// Nothing to delete
	if !ok {
		return false
	}

	delete(m.store, key)
	return true
}

// Check if a cache entry is expired. Note that this method is time-based and thus inherently
// stateful.
func (e *cacheEntry) isExpired() bool {
	if e.expiry.Unix() == 0 {
		return false
	}

	return e.expiry.Before(time.Now())
}
