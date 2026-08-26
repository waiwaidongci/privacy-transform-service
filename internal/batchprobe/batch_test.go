package batchprobe

import (
	"context"
	"testing"

	"github.com/ali/go-0821/privacy-transform-service/internal/adapter/cache"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

func TestMemoryCacheGetNilReceiver(t *testing.T) {
	var m *cache.Memory
	if v, ok := m.Get(context.Background(), "k"); ok || v != nil {
		t.Fatal("nil receiver Get should return not-ok and nil")
	}
}

func TestMemoryCacheSetNilReceiver(t *testing.T) {
	var m *cache.Memory
	m.Set(context.Background(), "k", domain.TransformRevision{Number: 1, Value: "x"})
}

func TestMemoryCacheDeleteNilReceiver(t *testing.T) {
	var m *cache.Memory
	m.Delete(context.Background(), "k")
}

func TestTTLCacheGetNilReceiver(t *testing.T) {
	var c *cache.TTL
	if v, ok := c.Get(context.Background(), "k"); ok || v != nil {
		t.Fatal("nil receiver Get should return not-ok and nil")
	}
}

func TestTTLCacheSetNilReceiver(t *testing.T) {
	var c *cache.TTL
	c.Set(context.Background(), "k", domain.TransformRevision{Number: 1, Value: "x"})
}

func TestTTLCacheDeleteNilReceiver(t *testing.T) {
	var c *cache.TTL
	c.Delete(context.Background(), "k")
}

func TestTTLCacheSizeNilReceiver(t *testing.T) {
	var c *cache.TTL
	if n := c.Size(); n != 0 {
		t.Fatalf("nil receiver Size should return 0, got %d", n)
	}
}
