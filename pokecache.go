package main

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	data     map[string]CacheEntry
	interval time.Duration
	mu       sync.RWMutex
	stopCh   chan struct{}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		data:     make(map[string]CacheEntry),
		interval: interval,
	}
	go cache.reapLoop()
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = CacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	data, exists := c.data[key]
	if !exists {
		return nil, false
	}
	return data.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.reap()
		case <-c.stopCh:
			return
		}
	}
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, entry := range c.data {
		if now.Sub(entry.createdAt) > c.interval {
			delete(c.data, key)
		}
	}
}

func (c *Cache) Stop() {
	close(c.stopCh)
}
