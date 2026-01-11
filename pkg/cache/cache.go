package cache

import (
	"sync"
)

type CacheInf interface {
	Get(key string) (any, bool)
	Set(key string, value any)
}

type cache struct {
	mu   sync.RWMutex
	data map[string]any
}

var Cache CacheInf = NewCache()

func NewCache() *cache {
	return &cache{
		data: make(map[string]any),
	}
}

func (c *cache) Get(name string) (country any, hasValue bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	country, hasValue = c.data[name]
	return
}

func (c *cache) Set(name string, country any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[name] = country
}
