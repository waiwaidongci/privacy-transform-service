package errorchainprobe

import (
	"context"
	"errors"
	"testing"

	"github.com/ali/go-0821/privacy-transform-service/internal/adapter/cache"
	"github.com/ali/go-0821/privacy-transform-service/internal/application"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/repository"
)

func newService() *application.Service {
	return application.NewService(repository.NewMemory(), cache.NewMemory())
}

func TestGetPolicyWorkspaceNotFoundChain(t *testing.T) {
	_, err := repository.NewMemory().GetPolicyWorkspace(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, err=%v", err)
	}
}

func TestGetTransformRuleSetNotFoundChain(t *testing.T) {
	_, err := repository.NewMemory().GetTransformRuleSet(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, err=%v", err)
	}
}

func TestGetTransformRevisionNotFoundChain(t *testing.T) {
	_, err := repository.NewMemory().GetTransformRevision(context.Background(), "missing", 1)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, err=%v", err)
	}
}

func TestCreatePolicyWorkspaceInvalidChain(t *testing.T) {
	err := newService().CreatePolicyWorkspace(context.Background(), domain.PolicyWorkspace{ID: "", Name: ""})
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestCreateRevisionNotFoundChain(t *testing.T) {
	_, err := newService().CreateTransformRevision(context.Background(), "missing", domain.TransformRevision{Value: "x"})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, err=%v", err)
	}
}

func TestEvaluateNotFoundChain(t *testing.T) {
	_, err := newService().Evaluate(context.Background(), "missing", domain.EvaluationContext{})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, err=%v", err)
	}
}

func seedRuleSet(t *testing.T) (*application.Service, domain.TransformRuleSet) {
	t.Helper()
	s := newService()
	rs := domain.TransformRuleSet{ID: "rs", PolicyWorkspaceID: "ws", ProcessingPurposeID: "p", Key: "k", Type: domain.TypeString, DefaultValue: "default"}
	if err := s.CreateTransformRuleSet(context.Background(), rs); err != nil {
		t.Fatalf("seed ruleset: %v", err)
	}
	return s, rs
}

func TestPublicationRulesetNotFoundChain(t *testing.T) {
	_, err := newService().PolicyPublication(context.Background(), "missing", 1, "p", "reason")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, err=%v", err)
	}
}

func TestPublicationRevisionNotFoundChain(t *testing.T) {
	s, _ := seedRuleSet(t)
	_, err := s.PolicyPublication(context.Background(), "rs", 99, "p", "reason")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, err=%v", err)
	}
}
