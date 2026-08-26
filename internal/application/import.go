// Package implementation for privacy transformation and sensitive-value protection.
package application

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

func (s *Service) ImportPolicyWorkspace(ctx context.Context, data []byte) error {
	var in Export
	if err := json.Unmarshal(data, &in); err != nil {
		return fmt.Errorf("%w: invalid import payload", domain.ErrInvalid)
	}
	if err := s.CreatePolicyWorkspace(ctx, in.PolicyWorkspace); err != nil {
		return fmt.Errorf("create workspace: %w", err)
	}
	for _, e := range in.ProcessingPurposes {
		if err := s.CreateProcessingPurpose(ctx, e); err != nil {
			return fmt.Errorf("create purpose: %w", err)
		}
	}
	for _, f := range in.TransformRuleSets {
		if err := s.CreateTransformRuleSet(ctx, f); err != nil {
			return fmt.Errorf("create ruleset: %w", err)
		}
	}
	return nil
}
func DecodeTransformRuleSet(data []byte) (domain.TransformRuleSet, error) {
	var f domain.TransformRuleSet
	err := json.Unmarshal(data, &f)
	return f, err
}
