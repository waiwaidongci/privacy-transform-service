package application

import (
	"context"
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"testing"
)

func TestTransformMasksHashesAndDeletes(t *testing.T) {
	s := NewPrivacyService()
	ctx := context.Background()
	p, err := s.PutPolicy(ctx, domain.TransformPolicy{ID: "p", Name: "PII", Rules: []domain.FieldRule{{Path: "email", Action: "mask", Preserve: 2}, {Path: "phone", Action: "hash", Salt: "v1"}, {Path: "secret", Action: "delete"}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err = s.PublishPolicy(p.ID); err != nil {
		t.Fatal(err)
	}
	got, err := s.Transform(ctx, "r1", p.ID, map[string]any{"email": "alice@example.com", "phone": "13800138000", "secret": "raw"})
	if err != nil {
		t.Fatal(err)
	}
	if got.Data["email"] == "alice@example.com" || got.Data["phone"] == "13800138000" {
		t.Fatal("sensitive values were not transformed")
	}
	if _, ok := got.Data["secret"]; ok {
		t.Fatal("deleted field remains")
	}
	if got.Summary.Transformed != 3 || got.Summary.Hashes != 1 || got.Summary.Deleted != 1 {
		t.Fatalf("bad summary: %#v", got.Summary)
	}
}
func TestDraftRejectedAndNestedArray(t *testing.T) {
	s := NewPrivacyService()
	ctx := context.Background()
	p, err := s.PutPolicy(ctx, domain.TransformPolicy{ID: "nested", Name: "nested", Rules: []domain.FieldRule{{Path: "users[*].email", Action: "mask", Preserve: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	input := map[string]any{"users": []any{map[string]any{"email": "alice@example.com"}, map[string]any{"email": "bob@example.com"}}}
	if _, err = s.Transform(ctx, "draft", p.ID, input); err == nil {
		t.Fatal("draft policy executed")
	}
	s.PublishPolicy(p.ID)
	got, err := s.Transform(ctx, "published", p.ID, input)
	if err != nil {
		t.Fatal(err)
	}
	users := got.Data["users"].([]any)
	if users[0].(map[string]any)["email"] == "alice@example.com" || users[1].(map[string]any)["email"] == "bob@example.com" {
		t.Fatalf("nested values unchanged: %#v", got.Data)
	}
	before := len(s.ListResults())
	if _, err = s.Simulate(ctx, p.ID, input); err != nil {
		t.Fatal(err)
	}
	if len(s.ListResults()) != before {
		t.Fatal("simulation persisted a processing record")
	}
}
func TestDeterministicToken(t *testing.T) {
	s := NewPrivacyService()
	ctx := context.Background()
	p, _ := s.PutPolicy(ctx, domain.TransformPolicy{ID: "p", Name: "tokens", Rules: []domain.FieldRule{{Path: "id", Action: "tokenize", Salt: "k"}}})
	s.PublishPolicy(p.ID)
	a, _ := s.Transform(ctx, "a", p.ID, map[string]any{"id": "42"})
	b, _ := s.Transform(ctx, "b", p.ID, map[string]any{"id": "42"})
	if a.Data["id"] != b.Data["id"] {
		t.Fatal("tokenization is not deterministic")
	}
}
