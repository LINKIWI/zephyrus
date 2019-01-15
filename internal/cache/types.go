package cache

import (
	"time"
)

// TTLCache formalizes an interface for a key-value cache backend that provides support for per-key
// time-based expiry (TTL).
type TTLCache interface {
	// Get retrieves the value associated with a key, if non-expired.
	Get(key string) interface{}

	// Set adds a new key-value pair with the specified TTL.
	Set(key string, value interface{}, ttl time.Duration)

	// Delete invalidates a cache entry by key.
	// Returns true if an entry was deleted; false otherwise.
	Delete(key string) bool
}
