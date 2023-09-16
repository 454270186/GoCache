package gocache

import (
	"sync"

	"github.com/454270186/GoCache/gocache/lru"
)

/*
	并发安全的cache
*/

type cache struct {
	mu         sync.Mutex
	lruCache   *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key, val string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.lruCache == nil {
		c.lruCache = lru.New(c.cacheBytes, nil)
	}
	c.lruCache.Put(key, val)
}

func (c *cache) get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if v, ok := c.lruCache.Get(key); ok {
		return v, true
	}

	return "", false
}
