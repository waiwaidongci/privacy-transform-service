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
func CloneRules(in []Rule) []Rule {
	out := make([]Rule, len(in))
	for i, r := range in {
		out[i] = r
		if r.Tags != nil {
			out[i].Tags = map[string]string{}
			for k, v := range r.Tags {
				out[i].Tags[k] = v
			}
		}
	}
	return out
}
func CopyTransformRuleSet(in TransformRuleSet) TransformRuleSet {
	in.Rules = CloneRules(in.Rules)
	return in
}
