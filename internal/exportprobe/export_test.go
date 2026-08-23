package exportprobe

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/ali/go-0821/privacy-transform-service/internal/application"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

type store struct {
	purposeErr error
	ruleSetErr error
}

func (s *store) CreatePolicyWorkspace(context.Context, domain.PolicyWorkspace) error { return nil }
func (s *store) GetPolicyWorkspace(context.Context, string) (domain.PolicyWorkspace, error) { return domain.PolicyWorkspace{ID: "w", Name: "workspace"}, nil }
func (s *store) CreateProcessingPurpose(context.Context, domain.ProcessingPurpose) error { return nil }
func (s *store) ListProcessingPurposes(context.Context, string) ([]domain.ProcessingPurpose, error) { return []domain.ProcessingPurpose{{ID: "p", PolicyWorkspaceID: "w", Name: "purpose"}}, s.purposeErr }
func (s *store) CreateTransformRuleSet(context.Context, domain.TransformRuleSet) error { return nil }
func (s *store) GetTransformRuleSet(context.Context, string) (domain.TransformRuleSet, error) { return domain.TransformRuleSet{}, domain.ErrNotFound }
func (s *store) ListTransformRuleSets(context.Context, string, string) ([]domain.TransformRuleSet, error) { return nil, s.ruleSetErr }
func (s *store) SaveTransformRevision(context.Context, domain.TransformRevision) error { return nil }
func (s *store) GetTransformRevision(context.Context, string, int) (domain.TransformRevision, error) { return domain.TransformRevision{}, domain.ErrNotFound }
func (s *store) ListTransformRevisions(context.Context, string) ([]domain.TransformRevision, error) { return nil, nil }
func (s *store) SavePolicyPublication(context.Context, domain.PolicyPublication) error { return nil }
func (s *store) ListPolicyPublications(context.Context, string) ([]domain.PolicyPublication, error) { return nil, nil }

type cache struct{}
func (cache) Get(context.Context, string) (*domain.TransformRevision, bool) { return nil, false }
func (cache) Set(context.Context, string, domain.TransformRevision) {}
func (cache) Delete(context.Context, string) {}

func TestExportChainPreservesPurposeRepositoryError(t *testing.T) {
	want := errors.New("purpose backend unavailable")
	_, err := application.NewService(&store{purposeErr: want}, cache{}).ExportPolicyWorkspace(context.Background(), "w")
	if !errors.Is(err, want) { t.Fatalf("lost purpose error: %v", err) }
}

func TestExportChainPreservesRuleSetRepositoryError(t *testing.T) {
	want := errors.New("rule set backend unavailable")
	_, err := application.NewService(&store{ruleSetErr: want}, cache{}).ExportPolicyWorkspace(context.Background(), "w")
	if !errors.Is(err, want) { t.Fatalf("lost rule set error: %v", err) }
}

func TestImportChainMarksInvalidPayload(t *testing.T) {
	var syntax *json.SyntaxError
	err := application.NewService(&store{}, cache{}).ImportPolicyWorkspace(context.Background(), []byte(`{"workspace":`))
	if !errors.Is(err, domain.ErrInvalid) || !errors.As(err, &syntax) { t.Fatalf("broken import chain: %v", err) }
}

func TestDecodeChainMarksInvalidPayload(t *testing.T) {
	var syntax *json.SyntaxError
	_, err := application.DecodeTransformRuleSet([]byte(`{"id":`))
	if !errors.Is(err, domain.ErrInvalid) || !errors.As(err, &syntax) { t.Fatalf("broken decode chain: %v", err) }
}
