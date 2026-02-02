package cache

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache()

	c.Set("country", "India")

	val, ok := c.Get("country")

	assert.True(t, ok)
	assert.Equal(t, "India", val)
}

func TestCache_GetNotFound(t *testing.T) {
	c := NewCache()

	val, ok := c.Get("missing")

	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestCache_EmptyKey(t *testing.T) {
	c := NewCache()

	c.Set("", "value")

	val, ok := c.Get("")

	assert.False(t, ok)
	assert.Nil(t, val)
}

func TestCache_OverwriteValue(t *testing.T) {
	c := NewCache()

	c.Set("k1", "v1")
	c.Set("k1", "v2")

	val, ok := c.Get("k1")

	assert.True(t, ok)
	assert.Equal(t, "v2", val)
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := NewCache()

	var wg sync.WaitGroup
	workers := 100

	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func(i int) {
			defer wg.Done()
			key := "key"

			c.Set(key, i)
			val, ok := c.Get(key)

			assert.True(t, ok)
			assert.NotNil(t, val)
		}(i)
	}

	wg.Wait()
}

func BenchmarkCache_GetSet(b *testing.B) {
	c := NewCache()

	for i := 0; i < b.N; i++ {
		c.Set("k", i)
		c.Get("k")
	}
}

// package cache

// import "testing"

// func TestCache_SetAndGet(t *testing.T) {
// 	c := NewCache()

// 	c.Set("India", "INR")
// 	val, ok := c.Get("India")
// 	if !ok {
// 		t.Fatalf("expected key to exist")
// 	}

// 	if val != "INR" {
// 		t.Fatalf("expected INR, got %v", val)
// 	}
// }

// func TestCache_GetMissing(t *testing.T) {
// 	c := NewCache()

// 	val, ok := c.Get("USA")
// 	if ok {
// 		t.Fatalf("expected key to be missing, got %v", val)
// 	}
// }

// func TestCache_Override(t *testing.T) {
// 	c := NewCache()

// 	c.Set("India", "INR")

// 	c.Set("India", "NewDelhi")
// 	val, ok := c.Get("India")
// 	if !ok {
// 		t.Fatalf("expected key to exist")
// 	}

// 	if val != "NewDelhi" {
// 		t.Fatalf("expected NewDelhi, got %v", val)
// 	}
// }
