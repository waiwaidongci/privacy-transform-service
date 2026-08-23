// Package implementation for privacy transformation and sensitive-value protection.
package domain

import (
	"encoding/json"
	"fmt"
)

func EncodeValue(v any) ([]byte, error) {
	b, e := json.Marshal(v)
	if e != nil {
		return nil, fmt.Errorf("encode value: %w", e)
	}
	return b, nil
}
func DecodeValue(data []byte, t ValueType) (any, error) {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("decode value: %w", err)
	}
	if err := ValidateValue(t, v); err != nil {
		return nil, err
	}
	return v, nil
}

// deepCloneValue returns an independent copy of a JSON-compatible value.
// Scalars are returned as-is; maps and slices are recursively copied so the
// caller can mutate the result without affecting the original.
func deepCloneValue(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, val := range x {
			out[k] = deepCloneValue(val)
		}
		return out
	case []any:
		out := make([]any, len(x))
		for i, val := range x {
			out[i] = deepCloneValue(val)
		}
		return out
	default:
		return v
	}
}

// CloneRules returns an independent deep copy of the given rule slice.
// Every nested reference type — Tags map, Percentage and time pointers, and
// the Value field — is duplicated so neither the caller nor the repository
// hold shared references into the other's copy.
func CloneRules(in []Rule) []Rule {
	if in == nil {
		return nil
	}
	out := make([]Rule, len(in))
	for i, r := range in {
		out[i] = r
		if r.Tags != nil {
			out[i].Tags = make(map[string]string, len(r.Tags))
			for k, v := range r.Tags {
				out[i].Tags[k] = v
			}
		}
		if r.Percentage != nil {
			p := *r.Percentage
			out[i].Percentage = &p
		}
		if r.StartAt != nil {
			t := *r.StartAt
			out[i].StartAt = &t
		}
		if r.EndAt != nil {
			t := *r.EndAt
			out[i].EndAt = &t
		}
		out[i].Value = deepCloneValue(r.Value)
	}
	return out
}

// CopyTransformRuleSet returns an independent deep copy of the rule set,
// including its Rules slice and DefaultValue.
func CopyTransformRuleSet(in TransformRuleSet) TransformRuleSet {
	in.Rules = CloneRules(in.Rules)
	in.DefaultValue = deepCloneValue(in.DefaultValue)
	return in
}

// CopyTransformRevision returns an independent deep copy of the revision,
// including its Rules slice and Value.
func CopyTransformRevision(in TransformRevision) TransformRevision {
	in.Value = deepCloneValue(in.Value)
	in.Rules = CloneRules(in.Rules)
	if in.PublishedAt != nil {
		t := *in.PublishedAt
		in.PublishedAt = &t
	}
	return in
}
