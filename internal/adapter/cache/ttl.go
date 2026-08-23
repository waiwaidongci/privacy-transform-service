// Package implementation for privacy transformation and sensitive-value protection.
package cache

import (
	"context"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"sync"
	"time"
)

type entry struct {
	revision domain.TransformRevision
	expires  time.Time
}
type TTL struct {
	mu   sync.RWMutex
	data map[string]entry
	ttl  time.Duration
}

func NewTTL(ttl time.Duration) *TTL {
	if ttl <= 0 {
		ttl = 5 * time.Minute
	}
	return &TTL{data: map[string]entry{}, ttl: ttl}
}
func (c *TTL) Get(_ context.Context, key string) (*domain.TransformRevision, bool) {
	c.mu.RLock()
	e, ok := c.data[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(e.expires) {
		if ok {
			c.Delete(context.Background(), key)
		}
		return nil, false
	}
	// Return a deep copy so the caller cannot mutate the cached revision's
	// Rules slice or Value.
	v := domain.CopyTransformRevision(e.revision)
	return &v, true
}
func (c *TTL) Set(_ context.Context, key string, v domain.TransformRevision) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// Store a deep copy so later mutation of the caller's Rules slice or
	// Value cannot leak into the cache.
	c.data[key] = entry{revision: domain.CopyTransformRevision(v), expires: time.Now().Add(c.ttl)}
}
func (c *TTL) Delete(_ context.Context, key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
func (c *TTL) Size() int { c.mu.RLock(); defer c.mu.RUnlock(); return len(c.data) }
