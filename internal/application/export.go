// Package implementation for privacy transformation and sensitive-value protection.
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

type Export struct {
	PolicyWorkspace    domain.PolicyWorkspace     `json:"workspace"`
	ProcessingPurposes []domain.ProcessingPurpose `json:"purposes"`
	TransformRuleSets  []domain.TransformRuleSet  `json:"ruleSets"`
}

func (s *Service) ExportPolicyWorkspace(ctx context.Context, id string) (Export, error) {
	p, e := s.store.GetPolicyWorkspace(ctx, id)
	if e != nil {
		return Export{}, fmt.Errorf("export workspace %s: %w", id, e)
	}
	envs, _ := s.store.ListProcessingPurposes(ctx, id)
	ruleSets := []domain.TransformRuleSet{}
	for _, env := range envs {
		items, _ := s.store.ListTransformRuleSets(ctx, id, env.ID)
		ruleSets = append(ruleSets, items...)
	}
	return Export{PolicyWorkspace: p, ProcessingPurposes: envs, TransformRuleSets: ruleSets}, nil
}
func (s *Service) ExportJSON(ctx context.Context, id string) ([]byte, error) {
	v, e := s.ExportPolicyWorkspace(ctx, id)
	if e != nil {
		return nil, fmt.Errorf("export json %s: %w", id, e)
	}
	return json.MarshalIndent(v, "", "  ")
}
