// Package implementation for privacy transformation and sensitive-value protection.
package cache

import (
	"context"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"sync"
)

type Memory struct {
	mu   sync.RWMutex
	data map[string]domain.TransformRevision
}

func NewMemory() *Memory { return &Memory{data: map[string]domain.TransformRevision{}} }
func (m *Memory) Get(_ context.Context, k string) (*domain.TransformRevision, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.data[k]
	if !ok {
		return nil, false
	}
	return &v, true
}
func (m *Memory) Set(_ context.Context, k string, v domain.TransformRevision) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[k] = v
}
func (m *Memory) Delete(_ context.Context, k string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, k)
}
