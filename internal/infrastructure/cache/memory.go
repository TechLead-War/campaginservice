package cache

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      []byte
	ExpiresAt time.Time
	CreatedAt time.Time
}

type MemoryCache struct {
	store       sync.Map
	maxSize     int
	currentSize int
	mutex       sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	cache := &MemoryCache{
		maxSize: 1000, // Default max items
	}

	// Start background cleanup
	go cache.cleanupExpired()

	return cache
}

func NewMemoryCacheWithSize(maxSize int) *MemoryCache {
	cache := &MemoryCache{
		maxSize: maxSize,
	}

	// Start background cleanup
	go cache.cleanupExpired()

	return cache
}

func (mc *MemoryCache) Set(key string, value []byte, ttl time.Duration) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	// Check if we need to evict items
	if mc.currentSize >= mc.maxSize {
		mc.evictOldest()
	}

	expiresAt := time.Now().Add(ttl)

	// Check if key already exists
	if _, exists := mc.store.Load(key); !exists {
		mc.currentSize++
	}

	mc.store.Store(key, CacheItem{
		Data:      value,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	})
}

// Get returns the data and true if found and not expired, otherwise returns nil and false
func (mc *MemoryCache) Get(key string) ([]byte, bool) {
	item, ok := mc.store.Load(key)

	if !ok {
		return nil, false
	}

	cacheItem, ok := item.(CacheItem)

	if !ok {
		mc.store.Delete(key) //Clean up malformed entry
		return nil, false
	}

	if time.Now().After(cacheItem.ExpiresAt) {
		mc.store.Delete(key)
		return nil, false
	}

	return cacheItem.Data, true
}

// Delete removes item from the cache
func (mc *MemoryCache) Delete(key string) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	if _, exists := mc.store.Load(key); exists {
		mc.currentSize--
	}
	mc.store.Delete(key)
}

// Size returns the current number of items in cache
func (mc *MemoryCache) Size() int {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()
	return mc.currentSize
}

// Clear removes all items from cache
func (mc *MemoryCache) Clear() {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	mc.store = sync.Map{}
	mc.currentSize = 0
}

func (mc *MemoryCache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time
	var found bool

	mc.store.Range(func(key, value interface{}) bool {
		if item, ok := value.(CacheItem); ok {
			if !found || item.CreatedAt.Before(oldestTime) {
				oldestKey = key.(string)
				oldestTime = item.CreatedAt
				found = true
			}
		}
		return true
	})

	if found {
		mc.store.Delete(oldestKey)
		mc.currentSize--
	}
}

func (mc *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mc.mutex.Lock()
		mc.store.Range(func(key, value interface{}) bool {
			if item, ok := value.(CacheItem); ok {
				if time.Now().After(item.ExpiresAt) {
					mc.store.Delete(key)
					mc.currentSize--
				}
			}
			return true
		})
		mc.mutex.Unlock()
	}
}
