// Package implementation for privacy transformation and sensitive-value protection.
package domain

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
	"time"
)

type EvaluationContext struct {
	SubjectID string            `json:"subject_id,omitempty"`
	Tags      map[string]string `json:"tags,omitempty"`
	Now       time.Time         `json:"now,omitempty"`
}

func Evaluate(ruleSet TransformRuleSet, revision *TransformRevision, ctx EvaluationContext) (any, int, error) {
	value := cloneJSONValue(ruleSet.DefaultValue)
	rules := ruleSet.Rules
	revisionNo := ruleSet.ActiveTransformRevision
	if revision != nil {
		value, rules, revisionNo = cloneJSONValue(revision.Value), revision.Rules, revision.Number
	}
	if ctx.Now.IsZero() {
		ctx.Now = time.Now()
	}
	sorted := make([]Rule, len(rules))
	copy(sorted, rules)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].Priority < sorted[j].Priority })
	for _, rule := range sorted {
		if !matchesTags(rule.Tags, ctx.Tags) || !matchesTime(rule, ctx.Now) || !matchesPercentage(rule.Percentage, ctx.SubjectID) {
			continue
		}
		return cloneJSONValue(rule.Value), revisionNo, nil
	}
	return value, revisionNo, nil
}

// cloneJSONValue returns a deep copy of v for the JSON-like container types
// (map/slice) so callers cannot mutate the stored input through the returned
// reference. Scalar values are returned as-is since they are immutable in Go.
func cloneJSONValue(v any) any {
	switch x := v.(type) {
	case map[string]any:
		return DeepCopyJSON(x)
	case []any:
		o := make([]any, len(x))
		for i, v := range x {
			o[i] = cloneJSONValue(v)
		}
		return o
	default:
		return x
	}
}

func matchesTags(expected, actual map[string]string) bool {
	for k, v := range expected {
		if actual == nil || actual[k] != v {
			return false
		}
	}
	return true
}
func matchesTime(r Rule, now time.Time) bool {
	if r.StartAt != nil && now.Before(*r.StartAt) {
		return false
	}
	if r.EndAt != nil && now.After(*r.EndAt) {
		return false
	}
	return true
}
func matchesPercentage(p *int, subject string) bool {
	if p == nil {
		return true
	}
	if *p <= 0 {
		return false
	}
	if *p >= 100 {
		return true
	}
	sum := sha256.Sum256([]byte(subject))
	n := binary.BigEndian.Uint16(sum[:2]) % 100
	return int(n) < *p
}

func ValidateRules(t ValueType, rules []Rule) error {
	for _, r := range rules {
		if err := ValidateValue(t, r.Value); err != nil {
			return fmt.Errorf("rule %s: %w", r.ID, err)
		}
		if r.Percentage != nil && (*r.Percentage < 0 || *r.Percentage > 100) {
			return fmt.Errorf("%w: percentage out of range", ErrInvalid)
		}
	}
	return nil
}
