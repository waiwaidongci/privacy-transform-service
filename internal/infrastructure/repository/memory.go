// Package implementation for privacy transformation and sensitive-value protection.
package repository

import (
	"context"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"sync"
)

type Memory struct {
	mu           sync.RWMutex
	workspaces   map[string]domain.PolicyWorkspace
	envs         map[string]domain.ProcessingPurpose
	ruleSets     map[string]domain.TransformRuleSet
	revisions    map[string][]domain.TransformRevision
	publications map[string][]domain.PolicyPublication
}

func NewMemory() *Memory {
	return &Memory{workspaces: map[string]domain.PolicyWorkspace{}, envs: map[string]domain.ProcessingPurpose{}, ruleSets: map[string]domain.TransformRuleSet{}, revisions: map[string][]domain.TransformRevision{}, publications: map[string][]domain.PolicyPublication{}}
}
func (m *Memory) CreatePolicyWorkspace(_ context.Context, p domain.PolicyWorkspace) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.workspaces[p.ID]; ok {
		return domain.ErrConflict
	}
	// Store an independent copy so later mutation of the caller's value
	// cannot leak into the repository.
	m.workspaces[p.ID] = clonePolicyWorkspace(p)
	return nil
}
func (m *Memory) GetPolicyWorkspace(_ context.Context, id string) (domain.PolicyWorkspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.workspaces[id]
	if !ok {
		return domain.PolicyWorkspace{}, domain.ErrNotFound
	}
	// Return a copy so the caller cannot mutate the stored value.
	return clonePolicyWorkspace(p), nil
}
func (m *Memory) CreateProcessingPurpose(_ context.Context, e domain.ProcessingPurpose) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.envs[e.ID]; ok {
		return domain.ErrConflict
	}
	m.envs[e.ID] = cloneProcessingPurpose(e)
	return nil
}
func (m *Memory) ListProcessingPurposes(_ context.Context, p string) ([]domain.ProcessingPurpose, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]domain.ProcessingPurpose, 0)
	for _, e := range m.envs {
		if e.PolicyWorkspaceID == p {
			// Append a copy so the caller cannot mutate the stored value.
			out = append(out, cloneProcessingPurpose(e))
		}
	}
	return out, nil
}
func (m *Memory) CreateTransformRuleSet(_ context.Context, f domain.TransformRuleSet) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if old, ok := m.ruleSets[f.ID]; ok {
		f.CreatedAt = old.CreatedAt
	}
	// Store a deep copy so later mutation of the caller's Rules slice or
	// DefaultValue cannot leak into the repository.
	m.ruleSets[f.ID] = cloneTransformRuleSet(f)
	return nil
}
func (m *Memory) GetTransformRuleSet(_ context.Context, id string) (domain.TransformRuleSet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.ruleSets[id]
	if !ok {
		return domain.TransformRuleSet{}, domain.ErrNotFound
	}
	// Return a deep copy so the caller cannot mutate the stored Rules slice
	// or DefaultValue.
	return cloneTransformRuleSet(f), nil
}
func (m *Memory) ListTransformRuleSets(_ context.Context, p, e string) ([]domain.TransformRuleSet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]domain.TransformRuleSet, 0)
	for _, f := range m.ruleSets {
		if f.PolicyWorkspaceID == p && (e == "" || f.ProcessingPurposeID == e) {
			// Append a deep copy so the caller cannot mutate the stored
			// Rules slices or DefaultValue values.
			out = append(out, cloneTransformRuleSet(f))
		}
	}
	return out, nil
}
func (m *Memory) SaveTransformRevision(_ context.Context, v domain.TransformRevision) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Store a deep copy so later mutation of the caller's Rules slice or
	// Value cannot leak into the repository.
	stored := cloneTransformRevision(v)
	arr := m.revisions[v.TransformRuleSetID]
	for i, x := range arr {
		if x.Number == v.Number {
			arr[i] = stored
			m.revisions[v.TransformRuleSetID] = arr
			return nil
		}
	}
	m.revisions[v.TransformRuleSetID] = append(arr, stored)
	return nil
}
func (m *Memory) GetTransformRevision(_ context.Context, id string, n int) (domain.TransformRevision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.revisions[id] {
		if v.Number == n {
			// Return a deep copy so the caller cannot mutate the stored
			// Rules slice or Value.
			return cloneTransformRevision(v), nil
		}
	}
	return domain.TransformRevision{}, domain.ErrNotFound
}
func (m *Memory) ListTransformRevisions(_ context.Context, id string) ([]domain.TransformRevision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	src := m.revisions[id]
	// Allocate a fresh slice and deep copy each entry so the caller cannot
	// mutate the stored revisions' Rules slices or Values.
	out := make([]domain.TransformRevision, len(src))
	for i, v := range src {
		out[i] = cloneTransformRevision(v)
	}
	return out, nil
}
func (m *Memory) SavePolicyPublication(_ context.Context, r domain.PolicyPublication) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Store a copy so later mutation of the caller's value cannot leak
	// into the repository.
	m.publications[r.TransformRuleSetID] = append(m.publications[r.TransformRuleSetID], clonePolicyPublication(r))
	return nil
}
func (m *Memory) ListPolicyPublications(_ context.Context, id string) ([]domain.PolicyPublication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	src := m.publications[id]
	// Allocate a fresh slice and copy each entry so the caller cannot
	// mutate the stored publications.
	out := make([]domain.PolicyPublication, len(src))
	for i, r := range src {
		out[i] = clonePolicyPublication(r)
	}
	return out, nil
}
