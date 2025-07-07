package internal

import "time"

func nowUTC() time.Time {
	return time.Now().UTC()
}

type cacheItem struct {
	value     []byte
	expiresAt time.Time
}

type Cache struct {
	items map[string]cacheItem
}

func (cache *Cache) Set(key string, value []byte, expiresIn time.Duration) {
	expiresAt := nowUTC().Add(expiresIn)
	cache.items[key] = cacheItem{value, expiresAt}
}

// hasValid item is true if and only if the item is in the cache, and it has not
// expired yet
func (cache *Cache) Get(key string) (item cacheItem, hasValidItem bool) {
	item, hasItem := cache.items[key]
	if hasItem && item.expiresAt.After(nowUTC()) {
		hasValidItem = true
	}
	return item, hasValidItem
}
