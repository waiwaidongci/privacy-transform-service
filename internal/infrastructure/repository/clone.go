// Package implementation for privacy transformation and sensitive-value protection.
package repository

import "github.com/ali/go-0821/privacy-transform-service/internal/domain"

// clonePolicyWorkspace returns a copy of p. PolicyWorkspace has no nested
// reference types, so a value copy is sufficient.
func clonePolicyWorkspace(p domain.PolicyWorkspace) domain.PolicyWorkspace { return p }

// cloneProcessingPurpose returns a copy of e. ProcessingPurpose has no nested
// reference types, so a value copy is sufficient.
func cloneProcessingPurpose(e domain.ProcessingPurpose) domain.ProcessingPurpose { return e }

// cloneTransformRevision returns an independent deep copy of v so that
// neither the caller nor the repository share references into each other's
// Rules slice or Value.
func cloneTransformRevision(v domain.TransformRevision) domain.TransformRevision {
	return domain.CopyTransformRevision(v)
}

// clonePolicyPublication returns a copy of r. PolicyPublication has no nested
// reference types, so a value copy is sufficient.
func clonePolicyPublication(r domain.PolicyPublication) domain.PolicyPublication { return r }

// cloneTransformRuleSet returns an independent deep copy of f so that
// neither the caller nor the repository share references into each other's
// Rules slice or DefaultValue.
func cloneTransformRuleSet(f domain.TransformRuleSet) domain.TransformRuleSet {
	return domain.CopyTransformRuleSet(f)
}
