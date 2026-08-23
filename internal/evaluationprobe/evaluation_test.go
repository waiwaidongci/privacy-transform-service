package evaluationprobe

import (
	"reflect"
	"testing"

	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

func rules() []domain.Rule {
	return []domain.Rule{
		{ID: "late", Priority: 20, Tags: map[string]string{"region": "jp"}, Value: map[string]any{"mode": "late"}},
		{ID: "early", Priority: 10, Tags: map[string]string{"region": "jp"}, Value: map[string]any{"mode": "early"}},
	}
}

func TestEvaluateDoesNotReorderRuleSet(t *testing.T) {
	set := domain.TransformRuleSet{DefaultValue: "default", Rules: rules()}
	before := []string{set.Rules[0].ID, set.Rules[1].ID}
	_, _, _ = domain.Evaluate(set, nil, domain.EvaluationContext{Tags: map[string]string{"region": "jp"}})
	after := []string{set.Rules[0].ID, set.Rules[1].ID}
	if !reflect.DeepEqual(before, after) {
		t.Fatalf("rule set reordered: %v -> %v", before, after)
	}
}

func TestEvaluateDoesNotReorderRevision(t *testing.T) {
	rev := domain.TransformRevision{Number: 2, Rules: rules()}
	_, _, _ = domain.Evaluate(domain.TransformRuleSet{}, &rev, domain.EvaluationContext{Tags: map[string]string{"region": "jp"}})
	if rev.Rules[0].ID != "late" || rev.Rules[1].ID != "early" {
		t.Fatalf("revision reordered: %#v", rev.Rules)
	}
}

func TestEvaluateDetachesMatchedValue(t *testing.T) {
	set := domain.TransformRuleSet{Rules: rules()}
	got, _, _ := domain.Evaluate(set, nil, domain.EvaluationContext{Tags: map[string]string{"region": "jp"}})
	got.(map[string]any)["mode"] = "changed"
	if set.Rules[1].Value.(map[string]any)["mode"] != "early" {
		t.Fatalf("matched value alias escaped")
	}
}

func TestEvaluateDetachesDefaultValue(t *testing.T) {
	set := domain.TransformRuleSet{DefaultValue: map[string]any{"scope": []any{"email"}}}
	got, _, _ := domain.Evaluate(set, nil, domain.EvaluationContext{})
	got.(map[string]any)["scope"].([]any)[0] = "phone"
	if set.DefaultValue.(map[string]any)["scope"].([]any)[0] != "email" {
		t.Fatalf("default value alias escaped")
	}
}
