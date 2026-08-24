// Package implementation for privacy transformation and sensitive-value protection.
package metrics

import (
	"fmt"
	"sort"
	"sync"
)

type Registry struct {
	mu       sync.RWMutex
	counters map[string]*uint64
}

func NewRegistry() *Registry { return &Registry{counters: map[string]*uint64{}} }
func (r *Registry) Counter(name string) *uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	v := uint64(0)
	r.counters[name] = &v
	return &v
}
func (r *Registry) Inc(name string) { (*r.Counter(name))++ }
func (r *Registry) Text() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := ""
	names := make([]string, 0, len(r.counters))
	for name := range r.counters {
		names = append(names, name)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(names)))
	for _, name := range names {
		out += fmt.Sprintf("%s %d\n", name, *r.counters[name])
	}
	return out
}
