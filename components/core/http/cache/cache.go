package cache

import (
	"time"
)

// Cacher is an interface for a cache for arbitrary objects
type Cacher interface {
	// Get returns the cached object under the key. It returns the object, its remaining
	// TTL. The returned interface is nil if of the object can't be found or it expired.
	// The error is non-nil if the object should be returned but can't due to some
	// implementation error.
	Get(key string) (interface{}, time.Duration, error)

	// Put adds an object under key to the cache. You have to provide its size. Any object
	// that may be stored under the same key will be overwritten.
	// The size parameter is a hint for how big this object is, whatever units
	// have been chosen for the MaxSize and MaxFileSize.
	// If the whole cache is bigger than the allowed size, the least used objects
	// will be removed until the cache size is lower than the allowed size.
	// A non-nil error is returned if the object couldn't be added to the cache.
	Put(key string, o interface{}, size uint64) error

	// Purge deletes the object stored under the key. A non-nil error is returned if
	// the object exists but can't be removed.
	Delete(key string) error

	// PurgeAll empties the whole cache.
	Purge()

	// TTL returns the cache's default TTL
	TTL() time.Duration

	// IsExtensionCacheable returns whether a file extension (e.g. .html) is allowed to be cached.
	IsExtensionCacheable(extension string) bool

	// IsSizeCacheable returns whether if a size is allowed to be cached.
	IsSizeCacheable(size uint64) bool
}
