package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var defaultConfig = LRUConfig{
	TTL:             time.Hour,
	MaxSize:         128,
	MaxFileSize:     0,
	AllowExtensions: []string{".html", ".js", ".jpg"},
	BlockExtensions: []string{".m3u8"},
	Logger:          nil,
}

func getCache(t *testing.T) *lrucache {
	cache, err := NewLRUCache(defaultConfig)
	require.Equal(t, nil, err)

	return cache.(*lrucache)
}

func TestNew(t *testing.T) {
	_, err := NewLRUCache(LRUConfig{
		TTL:         time.Hour,
		MaxSize:     128,
		MaxFileSize: 129,
	})
	require.NotEqual(t, nil, err)

	_, err = NewLRUCache(LRUConfig{
		TTL:         time.Hour,
		MaxSize:     0,
		MaxFileSize: 129,
	})
	require.Equal(t, nil, err)

	_, err = NewLRUCache(LRUConfig{
		TTL:         time.Hour,
		MaxSize:     128,
		MaxFileSize: 127,
	})
	require.Equal(t, nil, err)
}

func TestPut(t *testing.T) {
	cache := getCache(t)

	data, _, err := cache.Get("foobar")
	require.Equal(t, nil, data)
	require.Equal(t, nil, err)

	err = cache.Put("foobar", "hello", 64)
	require.Equal(t, nil, err)

	data, _, err = cache.Get("foobar")
	require.Equal(t, "hello", data)
	require.Equal(t, nil, err)
}

func TestPutTooBig(t *testing.T) {
	cache := getCache(t)

	err := cache.Put("foobar", "hello", 129)
	require.NotEqual(t, nil, err)
}

func TestPurge(t *testing.T) {
	cache := getCache(t)

	data, _, err := cache.Get("foobar")
	require.Equal(t, nil, data)
	require.Equal(t, nil, err)

	err = cache.Put("foobar", "hello", 64)
	require.Equal(t, nil, err)

	data, _, err = cache.Get("foobar")
	require.NotEqual(t, nil, data)
	require.Equal(t, nil, err)

	err = cache.Delete("foobar")
	require.Equal(t, nil, err)

	data, _, err = cache.Get("foobar")
	require.Equal(t, nil, data)
	require.Equal(t, nil, err)
}

func TestLRU(t *testing.T) {
	cache := getCache(t)

	err := cache.Put("1", "hello", 32)
	require.Equal(t, nil, err)

	err = cache.Put("2", "hello", 32)
	require.Equal(t, nil, err)

	err = cache.Put("3", "hello", 32)
	require.Equal(t, nil, err)

	err = cache.Put("4", "hello", 32)
	require.Equal(t, nil, err)

	data, _, _ := cache.Get("1")
	require.NotEqual(t, nil, data)

	data, _, _ = cache.Get("2")
	require.NotEqual(t, nil, data)

	data, _, _ = cache.Get("3")
	require.NotEqual(t, nil, data)

	data, _, _ = cache.Get("4")
	require.NotEqual(t, nil, data)

	data, _, _ = cache.Get("5")
	require.Equal(t, nil, data)

	err = cache.Put("5", "hello", 32)
	require.Equal(t, nil, err)

	data, _, _ = cache.Get("1")
	require.Equal(t, nil, data)

	data, _, _ = cache.Get("2")
	require.NotEqual(t, nil, data)

	data, _, _ = cache.Get("3")
	require.NotEqual(t, nil, data)

	data, _, _ = cache.Get("4")
	require.NotEqual(t, nil, data)

	data, _, _ = cache.Get("5")
	require.NotEqual(t, nil, data)
}

func TestAllowExtension(t *testing.T) {
	cache := getCache(t)

	r := cache.IsExtensionCacheable(".html")
	require.Equal(t, true, r)

	r = cache.IsExtensionCacheable(".png")
	require.Equal(t, false, r)
}

func TestBlockExtension(t *testing.T) {
	cache := getCache(t)
	cache.allowExtensions = []string{}

	r := cache.IsExtensionCacheable(".html")
	require.Equal(t, true, r)

	r = cache.IsExtensionCacheable(".m3u8")
	require.Equal(t, false, r)
}

func TestSize(t *testing.T) {
	cache := getCache(t)

	r := cache.IsSizeCacheable(127)
	require.Equal(t, true, r)

	r = cache.IsSizeCacheable(128)
	require.Equal(t, true, r)

	r = cache.IsSizeCacheable(129)
	require.Equal(t, false, r)
}
