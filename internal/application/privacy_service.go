// Package implementation for privacy transformation and sensitive-value protection.
package application

import (
	"context"
	"fmt"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"strings"
	"sync"
	"time"
)

type PrivacyService struct {
	mu            sync.RWMutex
	classes       map[string]domain.Classification
	policies      map[string]domain.TransformPolicy
	results       map[string]domain.TransformResult
	tokens        map[string]string
	policyHistory map[string][]domain.TransformPolicy
}

func NewPrivacyService() *PrivacyService {
	return &PrivacyService{classes: map[string]domain.Classification{}, policies: map[string]domain.TransformPolicy{}, results: map[string]domain.TransformResult{}, tokens: map[string]string{}, policyHistory: map[string][]domain.TransformPolicy{}}
}
func (s *PrivacyService) PutClassification(ctx context.Context, c domain.Classification) (domain.Classification, error) {
	if err := ctx.Err(); err != nil {
		return c, err
	}
	if c.ID == "" || c.Name == "" {
		return c, fmt.Errorf("%w: classification identity required", domain.ErrInvalid)
	}
	s.mu.Lock()
	s.classes[c.ID] = c
	s.mu.Unlock()
	return c, nil
}
func (s *PrivacyService) ListClassifications() []domain.Classification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o := make([]domain.Classification, 0, len(s.classes))
	for _, c := range s.classes {
		o = append(o, c)
	}
	return o
}
func (s *PrivacyService) PutPolicy(ctx context.Context, p domain.TransformPolicy) (domain.TransformPolicy, error) {
	if err := ctx.Err(); err != nil {
		return p, err
	}
	if p.ID == "" || p.Name == "" {
		return p, fmt.Errorf("%w: policy identity required", domain.ErrInvalid)
	}
	for _, r := range p.Rules {
		if err := domain.ValidateRule(r); err != nil {
			return p, err
		}
	}
	if p.Version < 1 {
		p.Version = 1
	}
	p.Status = "draft"
	p.UpdatedAt = time.Now().UTC()
	s.mu.Lock()
	s.policies[p.ID] = p
	s.policyHistory[p.ID] = append(s.policyHistory[p.ID], p)
	s.mu.Unlock()
	return p, nil
}
func (s *PrivacyService) PublishPolicy(id string) (domain.TransformPolicy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.policies[id]
	if !ok || p.Status != "draft" {
		return p, domain.ErrNotFound
	}
	p.Status = "published"
	p.UpdatedAt = time.Now().UTC()
	s.policies[id] = p
	return p, nil
}
func (s *PrivacyService) RollbackPolicy(id string) (domain.TransformPolicy, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	history := s.policyHistory[id]
	if len(history) < 2 {
		return domain.TransformPolicy{}, fmt.Errorf("%w: no prior published policy", domain.ErrConflict)
	}
	p := history[len(history)-2]
	p.Status = "published"
	p.UpdatedAt = time.Now().UTC()
	s.policies[id] = p
	s.policyHistory[id] = append(history, p)
	return p, nil
}
func (s *PrivacyService) Simulate(ctx context.Context, policyID string, input map[string]any) (domain.TransformResult, error) {
	result, err := s.Transform(ctx, "simulate-"+fmt.Sprint(time.Now().UnixNano()), policyID, input)
	if err == nil {
		s.mu.Lock()
		delete(s.results, result.RequestID)
		s.mu.Unlock()
	}
	return result, err
}
func (s *PrivacyService) ListResults() []domain.TransformResult {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]domain.TransformResult, 0, len(s.results))
	for _, r := range s.results {
		out = append(out, r)
	}
	return out
}
func (s *PrivacyService) Transform(ctx context.Context, requestID, policyID string, input map[string]any) (domain.TransformResult, error) {
	if err := ctx.Err(); err != nil {
		return domain.TransformResult{}, err
	}
	s.mu.RLock()
	p, ok := s.policies[policyID]
	s.mu.RUnlock()
	if !ok || p.Status != "published" {
		return domain.TransformResult{}, domain.ErrNotFound
	}
	out := domain.DeepCopyJSON(input)
	sum := domain.Summary{Processed: len(input)}
	for _, r := range p.Rules {
		count, err := domain.TransformPath(out, r.Path, func(v any) (any, bool) {
			sv := fmt.Sprint(v)
			switch r.Action {
			case "mask":
				return domain.MaskValue(sv, r.Preserve), true
			case "hash":
				sum.Hashes++
				return domain.HashValue(sv, r.Salt), true
			case "tokenize":
				s.mu.Lock()
				token := s.tokens[r.Path+"|"+sv]
				if token == "" {
					token = fmt.Sprintf("tok_%x", []byte(domain.HashValue(sv, r.Salt))[:12])
					s.tokens[r.Path+"|"+sv] = token
				}
				s.mu.Unlock()
				sum.Tokens++
				return token, true
			case "truncate":
				n := r.Preserve
				if n <= 0 {
					n = 4
				}
				rs := []rune(sv)
				if len(rs) > n {
					return string(rs[:n]), true
				}
				return sv, true
			case "generalize":
				return strings.Split(sv, " ")[0], true
			case "delete":
				sum.Deleted++
				return nil, false
			}
			return v, true
		})
		if err != nil {
			return domain.TransformResult{}, fmt.Errorf("apply %s: %w", r.Path, err)
		}
		sum.Transformed += count
	}
	res := domain.TransformResult{RequestID: requestID, PolicyID: policyID, Data: out, Summary: sum, CreatedAt: time.Now().UTC()}
	s.mu.Lock()
	s.results[requestID] = res
	s.mu.Unlock()
	return res, nil
}
func deletePath(m map[string]any, path string) {
	parts := strings.Split(path, ".")
	cur := m
	for i, p := range parts {
		if i == len(parts)-1 {
			delete(cur, p)
			return
		}
		n, ok := cur[p].(map[string]any)
		if !ok {
			return
		}
		cur = n
	}
}
func (s *PrivacyService) GetResult(id string) (domain.TransformResult, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.results[id]
	if !ok {
		return r, domain.ErrNotFound
	}
	return r, nil
}

func (s *PrivacyService) RevokeToken(token string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for key, value := range s.tokens {
		if value == token {
			delete(s.tokens, key)
			return nil
		}
	}
	return domain.ErrNotFound
}
