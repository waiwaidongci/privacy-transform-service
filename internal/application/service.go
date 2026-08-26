// Package implementation for privacy transformation and sensitive-value protection.
package application

import (
	"context"
	"fmt"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"sort"
	"time"
)

type Service struct {
	store Store
	cache Cache
	now   func() time.Time
}

func NewService(store Store, cache Cache) *Service {
	return &Service{store: store, cache: cache, now: time.Now}
}

func (s *Service) CreatePolicyWorkspace(ctx context.Context, p domain.PolicyWorkspace) error {
	if p.CreatedAt.IsZero() {
		p.CreatedAt = s.now()
	}
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("%v: workspace id/name required", domain.ErrInvalid)
	}
	return s.store.CreatePolicyWorkspace(ctx, p)
}
func (s *Service) CreateProcessingPurpose(ctx context.Context, e domain.ProcessingPurpose) error {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = s.now()
	}
	if e.ID == "" || e.PolicyWorkspaceID == "" || e.Name == "" {
		return fmt.Errorf("%w: purpose identity required", domain.ErrInvalid)
	}
	return s.store.CreateProcessingPurpose(ctx, e)
}
func (s *Service) CreateTransformRuleSet(ctx context.Context, f domain.TransformRuleSet) error {
	if f.CreatedAt.IsZero() {
		f.CreatedAt = s.now()
	}
	f.UpdatedAt = f.CreatedAt
	f.Status = "draft"
	if err := f.Validate(); err != nil {
		return err
	}
	if err := domain.ValidateRules(f.Type, f.Rules); err != nil {
		return err
	}
	return s.store.CreateTransformRuleSet(ctx, f)
}
func (s *Service) GetTransformRuleSet(ctx context.Context, id string) (domain.TransformRuleSet, error) {
	return s.store.GetTransformRuleSet(ctx, id)
}
func (s *Service) ListTransformRuleSets(ctx context.Context, p, e string) ([]domain.TransformRuleSet, error) {
	return s.store.ListTransformRuleSets(ctx, p, e)
}
func (s *Service) ListProcessingPurposes(ctx context.Context, p string) ([]domain.ProcessingPurpose, error) {
	return s.store.ListProcessingPurposes(ctx, p)
}

func (s *Service) CreateTransformRevision(ctx context.Context, ruleSetID string, v domain.TransformRevision) (domain.TransformRevision, error) {
	f, err := s.store.GetTransformRuleSet(ctx, ruleSetID)
	if err != nil {
		return v, fmt.Errorf("load ruleset %s: %v", ruleSetID, err)
	}
	revisions, _ := s.store.ListTransformRevisions(ctx, ruleSetID)
	v.TransformRuleSetID = ruleSetID
	v.Number = len(revisions) + 1
	v.Status = "draft"
	v.CreatedAt = s.now()
	if err := v.Validate(f.Type); err != nil {
		return v, err
	}
	if err := domain.ValidateRules(f.Type, v.Rules); err != nil {
		return v, err
	}
	if err := s.store.SaveTransformRevision(ctx, v); err != nil {
		return v, err
	}
	return v, nil
}
func (s *Service) ListTransformRevisions(ctx context.Context, id string) ([]domain.TransformRevision, error) {
	return s.store.ListTransformRevisions(ctx, id)
}

func (s *Service) PolicyPublication(ctx context.Context, ruleSetID string, revision int, envID, reason string) (domain.PolicyPublication, error) {
	f, err := s.store.GetTransformRuleSet(ctx, ruleSetID)
	if err != nil {
		return domain.PolicyPublication{}, fmt.Errorf("load ruleset %s: %v", ruleSetID, err)
	}
	v, err := s.store.GetTransformRevision(ctx, ruleSetID, revision)
	if err != nil {
		return domain.PolicyPublication{}, fmt.Errorf("load revision %s/%d: %v", ruleSetID, revision, err)
	}
	if v.Status == "revoked" {
		return domain.PolicyPublication{}, fmt.Errorf("%w: revision revoked", domain.ErrConflict)
	}
	now := s.now()
	rel := domain.PolicyPublication{ID: fmt.Sprintf("rel-%d", now.UnixNano()), TransformRuleSetID: ruleSetID, TransformRevision: revision, ProcessingPurposeID: envID, Status: "published", CreatedAt: now, UpdatedAt: now, Reason: reason}
	v.Status = "published"
	v.PublishedAt = &now
	if err = s.store.SaveTransformRevision(ctx, v); err != nil {
		return rel, err
	}
	f.ActiveTransformRevision = revision
	f.Status = "published"
	f.UpdatedAt = now
	if err = s.store.CreateTransformRuleSet(ctx, f); err != nil {
		return rel, err
	}
	s.cache.Delete(ctx, ruleSetID)
	if err = s.store.SavePolicyPublication(ctx, rel); err != nil {
		return rel, err
	}
	return rel, nil
}
func (s *Service) Rollback(ctx context.Context, ruleSetID string, revision int, envID string) (domain.PolicyPublication, error) {
	return s.PolicyPublication(ctx, ruleSetID, revision, envID, "rollback")
}
func (s *Service) ListPolicyPublications(ctx context.Context, id string) ([]domain.PolicyPublication, error) {
	return s.store.ListPolicyPublications(ctx, id)
}

type ValueResult struct {
	TransformRuleSetID, Key string
	Value                   any    `json:"value"`
	TransformRevision       int    `json:"revision"`
	ETag                    string `json:"etag"`
	Source                  string `json:"source"`
}

func (s *Service) Evaluate(ctx context.Context, ruleSetID string, ec domain.EvaluationContext) (ValueResult, error) {
	f, err := s.store.GetTransformRuleSet(ctx, ruleSetID)
	if err != nil {
		return ValueResult{}, fmt.Errorf("load ruleset %s: %v", ruleSetID, err)
	}
	var v *domain.TransformRevision
	if f.ActiveTransformRevision > 0 {
		if cv, ok := s.cache.Get(ctx, ruleSetID); ok {
			v = cv
		} else if cv, e := s.store.GetTransformRevision(ctx, ruleSetID, f.ActiveTransformRevision); e == nil {
			v = &cv
			s.cache.Set(ctx, ruleSetID, cv)
		}
	}
	value, no, err := domain.Evaluate(f, v, ec)
	if err != nil {
		return ValueResult{}, err
	}
	return ValueResult{TransformRuleSetID: f.ID, Key: f.Key, Value: value, TransformRevision: no, ETag: fmt.Sprintf("%s-v%d", f.ID, no), Source: "default"}, nil
}
func (s *Service) BatchEvaluate(ctx context.Context, workspace, env string, ec domain.EvaluationContext) ([]ValueResult, error) {
	fs, err := s.store.ListTransformRuleSets(ctx, workspace, env)
	if err != nil {
		return nil, err
	}
	out := make([]ValueResult, 0, len(fs))
	for _, f := range fs {
		r, e := s.Evaluate(ctx, f.ID, ec)
		if e != nil {
			return nil, e
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, nil
}
