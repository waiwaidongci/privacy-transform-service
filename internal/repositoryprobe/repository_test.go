package repositoryprobe

import (
	"context"
	"testing"

	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"github.com/ali/go-0821/privacy-transform-service/internal/infrastructure/repository"
)

func rev(number int) domain.TransformRevision {
	pct := 30
	return domain.TransformRevision{Number: number, TransformRuleSetID: "set", Value: map[string]any{"fields": []any{"email"}}, Rules: []domain.Rule{{ID: "r", Tags: map[string]string{"region": "jp"}, Percentage: &pct, Value: map[string]any{"mode": "mask"}}}}
}

func mutate(v *domain.TransformRevision) {
	v.Rules[0].Tags["region"] = "us"
	*v.Rules[0].Percentage = 99
	v.Value.(map[string]any)["fields"].([]any)[0] = "phone"
}

func clean(t *testing.T, got domain.TransformRevision) {
	t.Helper()
	if got.Rules[0].Tags["region"] != "jp" || *got.Rules[0].Percentage != 30 || got.Value.(map[string]any)["fields"].([]any)[0] != "email" {
		t.Fatalf("repository snapshot polluted: %#v", got)
	}
}

func TestRepositoryDetachesSavedRevision(t *testing.T) {
	m := repository.NewMemory()
	v := rev(1)
	if err := m.SaveTransformRevision(context.Background(), v); err != nil { t.Fatal(err) }
	mutate(&v)
	got, _ := m.GetTransformRevision(context.Background(), "set", 1)
	clean(t, got)
}

func TestRepositoryDetachesUpdatedRevision(t *testing.T) {
	m := repository.NewMemory()
	_ = m.SaveTransformRevision(context.Background(), rev(1))
	v := rev(1)
	_ = m.SaveTransformRevision(context.Background(), v)
	mutate(&v)
	got, _ := m.GetTransformRevision(context.Background(), "set", 1)
	clean(t, got)
}

func TestRepositoryDetachesLoadedRevision(t *testing.T) {
	m := repository.NewMemory()
	_ = m.SaveTransformRevision(context.Background(), rev(1))
	got, _ := m.GetTransformRevision(context.Background(), "set", 1)
	mutate(&got)
	again, _ := m.GetTransformRevision(context.Background(), "set", 1)
	clean(t, again)
}

func TestRepositoryDetachesListedRevision(t *testing.T) {
	m := repository.NewMemory()
	_ = m.SaveTransformRevision(context.Background(), rev(1))
	list, _ := m.ListTransformRevisions(context.Background(), "set")
	mutate(&list[0])
	again, _ := m.ListTransformRevisions(context.Background(), "set")
	clean(t, again[0])
}
