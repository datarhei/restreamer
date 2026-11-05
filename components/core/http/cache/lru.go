package cache

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/datarhei/core/v16/log"
)

// LRUConfig is the configuration for a new LRU cache
type LRUConfig struct {
	TTL             time.Duration // For how long the object should stay in cache
	MaxSize         uint64        // Max. size of the cache, 0 for unlimited, bytes
	MaxFileSize     uint64        // Max. file size allowed to put in cache, 0 for unlimited, bytes
	AllowExtensions []string      // List of file extension allowed to cache, empty list for all files
	BlockExtensions []string      // List of file extensions not allowed to cache, empty list for none
	Logger          log.Logger
}

type lrucache struct {
	ttl             time.Duration
	maxSize         uint64
	maxFileSize     uint64
	allowExtensions []string
	blockExtensions []string
	objects         map[string]*list.Element
	list            *list.List
	size            uint64
	lock            sync.Mutex
	logger          log.Logger
}

type value struct {
	key      string
	obj      interface{}
	expireAt time.Time
	size     uint64
}

// NewLRUCache returns an implementation of the Cacher interface that implements a LRU cache.
func NewLRUCache(config LRUConfig) (Cacher, error) {
	if config.MaxSize != 0 && config.MaxFileSize > config.MaxSize {
		return nil, fmt.Errorf("the max cache size has to be bigger than the max file size")
	}

	cache := &lrucache{
		ttl:         config.TTL,
		maxSize:     config.MaxSize,
		maxFileSize: config.MaxFileSize,
		list:        list.New(),
		objects:     make(map[string]*list.Element),
		logger:      config.Logger,
	}

	if cache.logger == nil {
		cache.logger = log.New("")
	}

	cache.allowExtensions = make([]string, len(config.AllowExtensions))
	copy(cache.allowExtensions, config.AllowExtensions)

	cache.blockExtensions = make([]string, len(config.BlockExtensions))
	copy(cache.blockExtensions, config.BlockExtensions)

	return cache, nil
}

// createValue create a value type from a key and object. This will be used as
// value for list.Element.
func (c *lrucache) createValue(key string, o interface{}, expireAt time.Time, size uint64) *value {
	v := &value{
		key:      key,
		obj:      o,
		expireAt: expireAt,
		size:     size,
	}

	return v
}

func (c *lrucache) Get(key string) (interface{}, time.Duration, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if the object exists
	elm, ok := c.objects[key]
	if !ok {
		return nil, c.ttl, nil
	}

	// Move it to the front of the list
	c.list.MoveToFront(elm)

	// Calculate the expiry date
	expire := elm.Value.(*value).expireAt

	// If it is expired, remove it from the list
	// and return as if it wasn't in the cache
	if expire.Before(time.Now()) {
		c.removeElement(elm)
		delete(c.objects, key)

		return nil, c.ttl, nil
	}

	// Get the actual cached object
	o := elm.Value.(*value).obj

	return o, time.Until(expire), nil
}

func (c *lrucache) Put(key string, o interface{}, size uint64) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if the object fits the cache
	if !c.IsSizeCacheable(size) {
		return fmt.Errorf("object too big to cache")
	}

	expireAt := time.Now().Add(c.ttl)

	// Check if we already have an object with this key in order
	// to replace it. Otherwise, create a new object.
	if elm, ok := c.objects[key]; ok {
		c.list.MoveToFront(elm)

		c.size -= elm.Value.(*value).size
		elm.Value = c.createValue(key, o, expireAt, size)
		c.size += elm.Value.(*value).size
	} else {
		elm = c.list.PushFront(c.createValue(key, o, expireAt, size))

		c.objects[key] = elm
		c.size += elm.Value.(*value).size
	}

	c.logger.WithFields(log.Fields{
		"key":        key,
		"size_bytes": size,
	}).Debug().Log("Added key")

	// If the size of the cache is exceeded, remove all least used
	// objects from the cache until the cache size is in its limits.
	if c.maxSize > 0 {
		for c.size > c.maxSize {
			elm := c.list.Back()
			if elm == nil {
				break
			}

			key := elm.Value.(*value).key

			c.logger.WithFields(log.Fields{
				"key":        key,
				"size_bytes": elm.Value.(*value).size,
			}).Debug().Log("Evicting key")

			c.removeElement(elm)
			delete(c.objects, key)
		}
	}

	return nil
}

func (c *lrucache) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	// Check if the object is in the cache. If not, do nothing.
	elm, ok := c.objects[key]
	if !ok {
		return nil
	}

	c.logger.WithFields(log.Fields{
		"key":        key,
		"size_bytes": elm.Value.(*value).size,
	}).Debug().Log("Purging key")

	c.removeElement(elm)
	delete(c.objects, key)

	return nil
}

func (c *lrucache) Purge() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.list.Init()
	c.objects = make(map[string]*list.Element)

	c.logger.WithField("size_bytes", c.size).Debug().Log("Purged all keys")

	c.size = 0
}

func (c *lrucache) TTL() time.Duration {
	return c.ttl
}

func (c *lrucache) IsExtensionCacheable(extension string) bool {
	if len(c.allowExtensions) == 0 && len(c.blockExtensions) == 0 {
		return true
	}

	for _, e := range c.blockExtensions {
		if extension == e {
			return false
		}
	}

	if len(c.allowExtensions) == 0 {
		return true
	}

	for _, e := range c.allowExtensions {
		if extension == e {
			return true
		}
	}

	return false
}

func (c *lrucache) IsSizeCacheable(size uint64) bool {
	// If the cache has a maximum size, the object can't be bigger than this size
	if c.maxSize != 0 && size > c.maxSize {
		return false
	}

	// If the cache has an object size limit, the object can't be bigger than this size
	if c.maxFileSize != 0 && size > c.maxFileSize {
		return false
	}

	return true
}

// removeElement removes an element from the list.
func (c *lrucache) removeElement(elm *list.Element) {
	c.list.Remove(elm)

	v := elm.Value.(*value)

	c.size -= v.size
}
