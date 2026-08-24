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
	m.workspaces[p.ID] = p
	return nil
}
func (m *Memory) GetPolicyWorkspace(_ context.Context, id string) (domain.PolicyWorkspace, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.workspaces[id]
	if !ok {
		return domain.PolicyWorkspace{}, domain.ErrNotFound
	}
	return p, nil
}
func (m *Memory) CreateProcessingPurpose(_ context.Context, e domain.ProcessingPurpose) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.envs[e.ID]; ok {
		return domain.ErrConflict
	}
	e2 := e
	m.envs[e.ID] = e2
	return nil
}
func (m *Memory) ListProcessingPurposes(_ context.Context, p string) ([]domain.ProcessingPurpose, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.ProcessingPurpose{}
	for _, e := range m.envs {
		if e.PolicyWorkspaceID == p {
			out = append(out, e)
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
	m.ruleSets[f.ID] = f
	return nil
}
func (m *Memory) GetTransformRuleSet(_ context.Context, id string) (domain.TransformRuleSet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	f, ok := m.ruleSets[id]
	if !ok {
		return domain.TransformRuleSet{}, domain.ErrNotFound
	}
	return f, nil
}
func (m *Memory) ListTransformRuleSets(_ context.Context, p, e string) ([]domain.TransformRuleSet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := []domain.TransformRuleSet{}
	for _, f := range m.ruleSets {
		if f.PolicyWorkspaceID == p && (e == "" || f.ProcessingPurposeID == e) {
			out = append(out, f)
		}
	}
	return out, nil
}
func (m *Memory) SaveTransformRevision(_ context.Context, v domain.TransformRevision) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	arr := m.revisions[v.TransformRuleSetID]
	for i, x := range arr {
		if x.Number == v.Number {
			arr[i] = v
			m.revisions[v.TransformRuleSetID] = arr
			return nil
		}
	}
	m.revisions[v.TransformRuleSetID] = append(arr, v)
	return nil
}
func (m *Memory) GetTransformRevision(_ context.Context, id string, n int) (domain.TransformRevision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.revisions[id] {
		if v.Number == n {
			return v, nil
		}
	}
	return domain.TransformRevision{}, domain.ErrNotFound
}
func (m *Memory) ListTransformRevisions(_ context.Context, id string) ([]domain.TransformRevision, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.TransformRevision(nil), m.revisions[id]...), nil
}
func (m *Memory) SavePolicyPublication(_ context.Context, r domain.PolicyPublication) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.publications[r.TransformRuleSetID] = append(m.publications[r.TransformRuleSetID], r)
	return nil
}
func (m *Memory) ListPolicyPublications(_ context.Context, id string) ([]domain.PolicyPublication, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]domain.PolicyPublication(nil), m.publications[id]...), nil
}
