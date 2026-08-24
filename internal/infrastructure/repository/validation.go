// Package implementation for privacy transformation and sensitive-value protection.
package repository

import (
	"fmt"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

func validatePolicyWorkspace(p domain.PolicyWorkspace) error {
	if p.ID == "" || p.Name == "" {
		return fmt.Errorf("%w: invalid workspace", domain.ErrInvalid)
	}
	return nil
}
func validateProcessingPurpose(e domain.ProcessingPurpose) error {
	if e.ID == "" || e.PolicyWorkspaceID == "" || e.Name == "" {
		return fmt.Errorf("%w: invalid purpose", domain.ErrInvalid)
	}
	return nil
}
func validateTransformRuleSet(f domain.TransformRuleSet) error { return f.Validate() }
func validateTransformRevision(v domain.TransformRevision, t domain.ValueType) error {
	return v.Validate(t)
}
