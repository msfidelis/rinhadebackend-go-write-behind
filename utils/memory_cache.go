package utils

import "sync"

type MemoryCache struct {
	items map[string]interface{}
	mutex sync.RWMutex
}

var cacheInstance *MemoryCache
var onceMemoryCache sync.Once

func GetCacheInstance() *MemoryCache {
	onceMemoryCache.Do(func() {
		cacheInstance = &MemoryCache{
			items: make(map[string]interface{}),
		}
	})
	return cacheInstance
}

func (c *MemoryCache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.items[key] = value
}

func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, found := c.items[key]
	return value, found
}
