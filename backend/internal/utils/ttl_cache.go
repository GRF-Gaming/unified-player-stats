package utils

import (
	"sync"
	"time"
)

type CacheItem[T any] struct {
	Value      *T
	Expiration time.Time
}

func (c *CacheItem[T]) IsExpired() bool {
	return c.Expiration.Before(time.Now())
}

type TTLCache[K comparable, T any] struct {
	items map[K]*CacheItem[T]
	mu    sync.RWMutex
	ttl   time.Duration
}

func NewTTLCache[K comparable, T any](ttl time.Duration) *TTLCache[K, T] {
	return &TTLCache[K, T]{
		items: make(map[K]*CacheItem[T]),
		mu:    sync.RWMutex{},
		ttl:   ttl,
	}
}

func (c *TTLCache[K, T]) Put(key K, val *T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = &CacheItem[T]{
		Value:      val,
		Expiration: time.Now().Add(c.ttl),
	}
}

func (c *TTLCache[K, T]) Get(key K) (*T, bool) {
	c.mu.RLock()
	i, exists := c.items[key]
	c.mu.RUnlock() // Unlock early

	if !exists {
		return nil, false
	}

	if i.IsExpired() {
		// Item is expired, need a write lock to delete
		c.mu.Lock()
		defer c.mu.Unlock()
		// Re-check in case another goroutine already deleted it
		i, exists = c.items[key]
		if exists && i.IsExpired() {
			delete(c.items, key)
		}
		return nil, false // Return as expired
	}

	return i.Value, true
}
