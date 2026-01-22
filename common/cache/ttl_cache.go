package cache

import (
	"sync"
	"time"
)

type TTLCache[T any] struct {
	mu    sync.RWMutex
	items map[string]ttlItem[T]
}

type ttlItem[T any] struct {
	value     T
	expiresAt time.Time
}

func NewTTLCache[T any]() *TTLCache[T] {
	return &TTLCache[T]{
		items: make(map[string]ttlItem[T]),
	}
}

func (c *TTLCache[T]) Get(key string) (T, bool) {
	var zero T
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return zero, false
	}
	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		c.mu.Lock()
		delete(c.items, key)
		c.mu.Unlock()
		return zero, false
	}
	return item.value, true
}

func (c *TTLCache[T]) Set(key string, value T, ttl time.Duration) {
	exp := time.Time{}
	if ttl > 0 {
		exp = time.Now().Add(ttl)
	}
	c.mu.Lock()
	c.items[key] = ttlItem[T]{value: value, expiresAt: exp}
	c.mu.Unlock()
}

