package internal

import (
	"time"

	"github.com/Jswig/gomemcached/internal/util"
)

type cacheItem struct {
	value     []byte
	expiresAt time.Time
}

type Cache struct {
	items map[string]cacheItem
}

func NewCache() *Cache {
	emptyItems := make(map[string]cacheItem)
	return &Cache{emptyItems}
}

func (cache *Cache) Set(key string, value []byte, expiresIn time.Duration) {
	expiresAt := util.NowUTC().Add(expiresIn)
	cache.items[key] = cacheItem{value, expiresAt}
}

// hasValid item is true if and only if the item is in the cache, and it has not
// expired yet
func (cache *Cache) Get(key string) (value []byte, isValidItem bool) {
	item, hasItem := cache.items[key]
	if hasItem && item.expiresAt.After(util.NowUTC()) {
		isValidItem = true
	}
	return item.value, isValidItem
}
