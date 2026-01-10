package cache

import (
	"sync"
)

type CacheInf interface {
	Get(key string) (any, bool)
	Set(key string, value any)
}

type Cache struct {
	mu        sync.RWMutex
	countries map[string]any
}

var CountryCache CacheInf = NewCache()

func NewCache() *Cache {
	return &Cache{
		countries: make(map[string]any),
	}
}

func (c *Cache) Get(name string) (country any, hasValue bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	country, ok := c.countries[name]
	return country, ok
}

func (c *Cache) Set(name string, country any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.countries[name] = country
}
