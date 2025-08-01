package cache

import (
	"context"
	"sync"
	"time"
)

type InMemoryCache struct {
	mu    sync.Mutex
	store map[string]interface{}
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		store: make(map[string]interface{}),
	}
}

func (c *InMemoryCache) Get(_ context.Context, key string) (interface{}, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.store[key]
	if !ok {
		return nil, ErrNotFound
	}
	return val, nil
}

func (c *InMemoryCache) Set(_ context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
	return nil
}

func (c *InMemoryCache) Delete(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
	return nil
}

func (c *InMemoryCache) Has(ctx context.Context, key string) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.store[key]
	return ok, nil
}
