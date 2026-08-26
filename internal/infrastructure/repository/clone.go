// Package implementation for privacy transformation and sensitive-value protection.
package repository

import "github.com/ali/go-0821/privacy-transform-service/internal/domain"

func clonePolicyWorkspace(p domain.PolicyWorkspace) domain.PolicyWorkspace       { return p }
func cloneProcessingPurpose(e domain.ProcessingPurpose) domain.ProcessingPurpose { return e }
func cloneTransformRevision(v domain.TransformRevision) domain.TransformRevision {
	v.Rules = domain.CloneRules(v.Rules)
	return v
}
func clonePolicyPublication(r domain.PolicyPublication) domain.PolicyPublication { return r }
func cloneTransformRuleSet(f domain.TransformRuleSet) domain.TransformRuleSet {
	return domain.CopyTransformRuleSet(f)
}
