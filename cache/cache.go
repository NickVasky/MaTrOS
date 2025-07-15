package cache

import (
	"sync"

	"github.com/emersion/go-imap/v2"
)

type InMemoryCache struct {
	mu    sync.Mutex
	store map[imap.UID]interface{}
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		store: make(map[imap.UID]interface{}),
	}
}

func (c *InMemoryCache) Get(key imap.UID) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.store[key]
	return val, ok
}

func (c *InMemoryCache) Set(key imap.UID, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[key] = value
}

func (c *InMemoryCache) Delete(key imap.UID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}

func (c *InMemoryCache) Has(key imap.UID) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	_, ok := c.store[key]
	return ok
}
