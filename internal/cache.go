package internal

import (
	"sync"
	"time"

	"github.com/Jswig/gomemcached/internal/util"
)

type cacheItem struct {
	value []byte
	// Zero time here represents an item that never expires
	expiresAt time.Time
}

func NeverExpires() time.Time {
	// Zero time represents an item that never expires
	return time.Time{}
}

type Cache struct {
	items map[string]cacheItem
	mu    sync.RWMutex
}

func NewCache() *Cache {
	emptyItems := make(map[string]cacheItem)
	return &Cache{items: emptyItems}
}

// checks if a cache item has expired
func isValid(item cacheItem) bool {
	// zero time means that the item should never expire
	if item.expiresAt.Equal(NeverExpires()) {
		return true
	} else {
		return item.expiresAt.After(util.NowUTC())
	}
}

// returns true if and only if the item already was not already in the cache and
// was added
func (cache *Cache) Add(key string, value []byte, expiresAt time.Time) (wasAdded bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	_, exists := cache.items[key]
	if !exists {
		cache.items[key] = cacheItem{value, expiresAt}
		wasAdded = true
	}
	return
}

// returns true if and only if the item already existed in cache
func (cache *Cache) Delete(key string) (wasDeleted bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	_, exists := cache.items[key]
	if exists {
		delete(cache.items, key)
		wasDeleted = true
	}
	return
}

// hasValidItem is true if and only if the item is in the cache, and it has not
// expired yet
func (cache *Cache) Get(key string) (value []byte, hasValidItem bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	item, exists := cache.items[key]
	if exists && isValid(item) {
		hasValidItem = true
		value = item.value
	}
	return value, hasValidItem
}

// wasReplaced is true if and only if the item was already in the cache and was
// replaced
func (cache *Cache) Replace(
	key string, value []byte, expiresAt time.Time,
) (wasReplaced bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	_, exists := cache.items[key]
	if exists {
		cache.items[key] = cacheItem{value, expiresAt}
		wasReplaced = true
	}
	return
}

func (cache *Cache) Set(key string, value []byte, expiresAt time.Time) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	cache.items[key] = cacheItem{value, expiresAt}
}
