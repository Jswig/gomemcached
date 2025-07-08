package internal

import (
	"time"
	"sync"
	"github.com/Jswig/gomemcached/internal/util"
)

type cacheItem struct {
	value     []byte
	expiresAt time.Time
}

type Cache struct {
	items map[string]cacheItem
	mu sync.RWMutex
}

func NewCache() *Cache {
	emptyItems := make(map[string]cacheItem)
	return &Cache{items: emptyItems}
}

func (cache *Cache) Add(key string, value []byte, expiresIn time.Duration) (wasAdded bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	expiresAt := util.NowUTC().Add(expiresIn)
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

// hasValid item is true if and only if the item is in the cache, and it has not
// expired yet
func (cache *Cache) Get(key string) (value []byte, isValidItem bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()
	item, exists := cache.items[key]
	if exists && item.expiresAt.After(util.NowUTC()) {
		isValidItem = true
	}
	return item.value, isValidItem
}

func (cache *Cache) GetAndTouch(key string, expiresIn time.Duration) (value []byte, isValidItem bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	item, exists := cache.items[key]
	if exists && item.expiresAt.After(util.NowUTC()) {
		isValidItem = true
		expiresAt := util.NowUTC().Add(expiresIn)
		cache.items[key] = cacheItem{item.value, expiresAt}
		value = item.value
	}
	return value, isValidItem
}


func (cache *Cache) Replace(key string, value []byte, expiresIn time.Duration) (wasReplaced bool) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	expiresAt := util.NowUTC().Add(expiresIn)
	_, exists := cache.items[key]
	if exists {
		cache.items[key] = cacheItem{value, expiresAt}
		wasReplaced = true
	}
	return
}

func (cache *Cache) Set(key string, value []byte, expiresIn time.Duration) {
	cache.mu.Lock()
	defer cache.mu.Unlock()
	expiresAt := util.NowUTC().Add(expiresIn)
	cache.items[key] = cacheItem{value, expiresAt}
}
