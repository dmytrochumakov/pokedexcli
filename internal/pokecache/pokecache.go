package pokecache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

// Cache -
type Cache struct {
	data map[string]CacheEntry
	mu   *sync.Mutex
}

// NewCache -
func NewCache(interval time.Duration) Cache {
	cache := Cache{
		data: make(map[string]CacheEntry),
		mu:   &sync.Mutex{},
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheEntry{
		createdAt: time.Now().UTC(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	data, exists := c.data[key]
	return data.val, exists
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(now time.Time, last time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.data {
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.data, k)
		}
	}
}
