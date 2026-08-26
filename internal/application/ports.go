// Package implementation for privacy transformation and sensitive-value protection.
package application

import (
	"context"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

type Store interface {
	CreatePolicyWorkspace(context.Context, domain.PolicyWorkspace) error
	GetPolicyWorkspace(context.Context, string) (domain.PolicyWorkspace, error)
	CreateProcessingPurpose(context.Context, domain.ProcessingPurpose) error
	ListProcessingPurposes(context.Context, string) ([]domain.ProcessingPurpose, error)
	CreateTransformRuleSet(context.Context, domain.TransformRuleSet) error
	GetTransformRuleSet(context.Context, string) (domain.TransformRuleSet, error)
	ListTransformRuleSets(context.Context, string, string) ([]domain.TransformRuleSet, error)
	SaveTransformRevision(context.Context, domain.TransformRevision) error
	GetTransformRevision(context.Context, string, int) (domain.TransformRevision, error)
	ListTransformRevisions(context.Context, string) ([]domain.TransformRevision, error)
	SavePolicyPublication(context.Context, domain.PolicyPublication) error
	ListPolicyPublications(context.Context, string) ([]domain.PolicyPublication, error)
}

type Cache interface {
	Get(context.Context, string) (*domain.TransformRevision, bool)
	Set(context.Context, string, domain.TransformRevision)
	Delete(context.Context, string)
}
