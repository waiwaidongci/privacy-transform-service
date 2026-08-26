package exportimportprobe

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

func TestExportNotFoundKeepsChain(t *testing.T) {
	_, err := newService().ExportPolicyWorkspace(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, err=%v", err)
	}
}

func TestExportJSONNotFoundKeepsChain(t *testing.T) {
	_, err := newService().ExportJSON(context.Background(), "missing")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("errors.Is(err, ErrNotFound) = false, err=%v", err)
	}
}

func TestImportInvalidPayloadKeepsChain(t *testing.T) {
	err := newService().ImportPolicyWorkspace(context.Background(), []byte("not-json"))
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestImportInvalidWorkspaceKeepsChain(t *testing.T) {
	err := newService().ImportPolicyWorkspace(context.Background(), []byte(`{}`))
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestImportInvalidPurposeKeepsChain(t *testing.T) {
	payload := []byte(`{"workspace":{"id":"w","name":"n"},"purposes":[{"id":"","workspace_id":"w","name":""}]}`)
	err := newService().ImportPolicyWorkspace(context.Background(), payload)
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestImportInvalidRuleSetKeepsChain(t *testing.T) {
	payload := []byte(`{"workspace":{"id":"w","name":"n"},"purposes":[{"id":"p","workspace_id":"w","name":"pn"}],"ruleSets":[{"id":"","workspace_id":"w","purpose_id":"p","key":"k","type":"string","default_value":"x"}]}`)
	err := newService().ImportPolicyWorkspace(context.Background(), payload)
	if !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}
