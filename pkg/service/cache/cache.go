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

func NewCache() CacheInf {
	return &cache{
		data: make(map[string]any),
	}
}

func (c *cache) Get(key string) (country any, hasValue bool) {
	if key == "" {
		return nil, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()

	country, hasValue = c.data[key]
	return
}

func (c *cache) Set(key string, value any) {
	if key == "" {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value
}
