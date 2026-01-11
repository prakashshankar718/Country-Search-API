package cache

import "testing"

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache()

	c.Set("India", "INR")
	val, ok := c.Get("India")
	if !ok {
		t.Fatalf("expected key to exist")
	}

	if val != "INR" {
		t.Fatalf("expected INR, got %v", val)
	}
}

func TestCache_GetMissing(t *testing.T) {
	c := NewCache()

	val, ok := c.Get("USA")
	if ok {
		t.Fatalf("expected key to be missing, got %v", val)
	}
}

func TestCache_Override(t *testing.T) {
	c := NewCache()

	c.Set("India", "INR")

	c.Set("India", "NewDelhi")
	val, ok := c.Get("India")
	if !ok {
		t.Fatalf("expected key to exist")
	}

	if val != "NewDelhi" {
		t.Fatalf("expected NewDelhi, got %v", val)
	}
}
