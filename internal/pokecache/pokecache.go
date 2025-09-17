package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	Entries map[string]cacheEntry
	mu      sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		Entries: make(map[string]cacheEntry),
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			c.reapLoop(interval)
		}
	}()
	return c
}

func (c *Cache) reapLoop(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	for key, entry := range c.Entries {
		if now.Sub(entry.createdAt) > interval {
			delete(c.Entries, key)
		}
	}
}

func (c *Cache) Add(key string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if value, ok := c.Entries[key]; ok {
		return value.val, true
	} else {
		return nil, false
	}
}
