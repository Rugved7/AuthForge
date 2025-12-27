package cache

import (
	"context"
	"sync"
	"time"
)

type item struct {
	value     string
	expiresAt time.Time
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]item
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]item),
	}
}

func (c *MemoryCache) Get(ctx context.Context, key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.RUnlock()

	it, ok := c.items[key]
	if !ok || time.Now().After(it.expiresAt) {
		return "", false
	}
	return it.value, true
}

func (c *MemoryCache) Set(ctx context.Context, key, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
}

func (c *MemoryCache) Delete(ctx context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}
