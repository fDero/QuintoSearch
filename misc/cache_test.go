package misc

import (
	"testing"
)

func TestCreateEmptyLRUCache(t *testing.T) {
	cache := NewLRUCache[int](300)
	if cache.maxItems != 300 {
		t.Errorf("Expected maxItems to be 300, got %d", cache.maxItems)
	}
	if cache.mru != nil {
		t.Error("Expected mru to be nil")
	}
	if cache.lru != nil {
		t.Error("Expected lru to be nil")
	}
	if len(cache.storage) != 0 {
		t.Errorf("Expected storage to be empty, got %d", len(cache.storage))
	}
}

func TestInsertOneElementInLRUCache(t *testing.T) {
	cache := NewLRUCache[int](300)
	cache.Store("key", 42)
	if cache.mru == nil {
		t.Error("Expected mru to be set")
	}
	if cache.lru == nil {
		t.Error("Expected lru to be set")
	}
	if len(cache.storage) != 1 {
		t.Errorf("Expected storage to have 1 element, got %d", len(cache.storage))
	}
	if cache.mru != cache.lru {
		t.Error("Expected mru and lru to be the same")
	}
	if _, exists := cache.storage["key"]; !exists {
		t.Error("Expected key to be in the cache")
	}
}

func TestInsertMultipleElementsInLRUCache(t *testing.T) {
	cache := NewLRUCache[int](300)
	cache.Store("key1", 42)
	cache.Store("key2", 43)
	if cache.mru == nil {
		t.Error("Expected mru to be set")
	}
	if cache.lru == nil {
		t.Error("Expected lru to be set")
	}
	if len(cache.storage) != 2 {
		t.Errorf("Expected storage to have 2 elements, got %d", len(cache.storage))
	}
	if cache.mru == cache.lru {
		t.Error("Expected mru and lru to be different")
	}
	if _, exists := cache.storage["key1"]; !exists {
		t.Error("Expected key1 to be in the cache")
	}
	if _, exists := cache.storage["key2"]; !exists {
		t.Error("Expected key2 to be in the cache")
	}
}

func TestInsertMultipleElementsInLRUCache2(t *testing.T) {
	cache := NewLRUCache[int](2)
	cache.Store("key1", 41)
	cache.Store("key2", 42)
	cache.Store("key3", 43)
	if cache.mru == nil {
		t.Error("Expected mru to be set")
	}
	if cache.lru == nil {
		t.Error("Expected lru to be set")
	}
	if len(cache.storage) != 2 {
		t.Errorf("Expected storage to have 2 elements, got %d", len(cache.storage))
	}
	if cache.mru == cache.lru {
		t.Error("Expected mru and lru to be different")
	}
	if _, exists := cache.storage["key1"]; exists {
		t.Error("Expected key1 to be evicted")
	}
	if _, exists := cache.storage["key2"]; !exists {
		t.Error("Expected key2 to be in the cache")
	}
	if _, exists := cache.storage["key3"]; !exists {
		t.Error("Expected key3 to be in the cache")
	}
}

func TestInsertRetrieveWorksOnLRUCache(t *testing.T) {
	cache := NewLRUCache[int](4)
	cache.Store("key1", 41)
	cache.Store("key2", 42)
	cache.Store("key3", 43)
	cache.Store("key4", 43)
	if cache.mru == nil || cache.lru == nil {
		t.Error("Expected mru/lru to be set")
	}
	if len(cache.storage) != 4 {
		t.Errorf("Expected storage to have 4 elements, got %d", len(cache.storage))
	}
	if cache.mru == cache.lru {
		t.Error("Expected mru and lru to be different")
	}
	if _, exists := cache.storage["key1"]; !exists {
		t.Error("Expected key1 to be in the cache")
	}
	if _, exists := cache.storage["key2"]; !exists {
		t.Error("Expected key2 to be in the cache")
	}
	if _, exists := cache.storage["key3"]; !exists {
		t.Error("Expected key3 to be in the cache")
	}
	if _, exists := cache.storage["key4"]; !exists {
		t.Error("Expected key4 to be in the cache")
	}
	if value := *cache.Retrieve("key1"); value != 41 {
		t.Errorf("Expected key1 to be 41, got %d", value)
	}
	if value := *cache.Retrieve("key2"); value != 42 {
		t.Errorf("Expected key2 to be 42, got %d", value)
	}
	if value := *cache.Retrieve("key3"); value != 43 {
		t.Errorf("Expected key3 to be 43, got %d", value)
	}
	if value := *cache.Retrieve("key4"); value != 43 {
		t.Errorf("Expected key4 to be 43, got %d", value)
	}
}

func TestRetrieveCausesLRUToBecomeMRU(t *testing.T) {
	cache := NewLRUCache[int](4)
	cache.Store("key1", 41)
	cache.Store("key2", 42)
	cache.Store("key3", 43)
	cache.Store("key4", 44)
	if cache.mru == nil {
		t.Error("Expected mru to be set")
	}
	if cache.mru == nil {
		t.Error("Expected mru to be set")
	}
	cache.Retrieve("key1")
	if cache.mru == nil {
		t.Error("Expected mru to be set")
	}
	if cache.lru == nil {
		t.Error("Expected lru to be set")
	}
	if cache.mru.item != 41 {
		t.Errorf("Expected mru to be 41, got %d", cache.mru.item)
	}
	if cache.lru.item != 42 {
		t.Errorf("Expected lru to be 42, got %d", cache.lru.item)
	}
}

func TestRetrieveLeavesLRUuntouched(t *testing.T) {
	cache := NewLRUCache[int](4)
	cache.Store("key1", 41)
	cache.Store("key2", 42)
	cache.Store("key3", 43)
	cache.Store("key4", 44)
	if cache.mru == nil {
		t.Error("Expected mru to be set")
	}
	if cache.mru == nil {
		t.Error("Expected mru to be set")
	}
	cache.Retrieve("key2")
	if cache.mru == nil {
		t.Error("Expected mru to be set")
	}
	if cache.lru == nil {
		t.Error("Expected lru to be set")
	}
	if cache.mru.item != 42 {
		t.Errorf("Expected mru to be 42, got %d", cache.mru.item)
	}
	if cache.lru.item != 41 {
		t.Errorf("Expected lru to be 41, got %d", cache.lru.item)
	}
}
