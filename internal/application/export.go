// Package implementation for privacy transformation and sensitive-value protection.
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

type Export struct {
	PolicyWorkspace    domain.PolicyWorkspace      `json:"workspace"`
	ProcessingPurposes []domain.ProcessingPurpose  `json:"purposes"`
	TransformRuleSets  []domain.TransformRuleSet   `json:"ruleSets"`
}

func (s *Service) ExportPolicyWorkspace(ctx context.Context, id string) (Export, error) {
	p, err := s.store.GetPolicyWorkspace(ctx, id)
	if err != nil {
		return Export{}, err
	}
	envs, err := s.store.ListProcessingPurposes(ctx, id)
	if err != nil {
		return Export{}, fmt.Errorf("export workspace %q: list purposes: %w", id, err)
	}
	ruleSets := []domain.TransformRuleSet{}
	for _, env := range envs {
		items, err := s.store.ListTransformRuleSets(ctx, id, env.ID)
		if err != nil {
			return Export{}, fmt.Errorf("export workspace %q: list rule sets for purpose %q: %w", id, env.ID, err)
		}
		ruleSets = append(ruleSets, items...)
	}
	return Export{PolicyWorkspace: p, ProcessingPurposes: envs, TransformRuleSets: ruleSets}, nil
}
func (s *Service) ExportJSON(ctx context.Context, id string) ([]byte, error) {
	v, err := s.ExportPolicyWorkspace(ctx, id)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(v, "", "  ")
}
