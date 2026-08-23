// Package implementation for privacy transformation and sensitive-value protection.
package domain

import (
	"sort"
	"strings"
)

type TransformRuleSetFilter struct {
	KeyContains         string
	Status              string
	Type                ValueType
	ProcessingPurposeID string
	SortBy              string
	Descending          bool
}

func FilterTransformRuleSets(ruleSets []TransformRuleSet, f TransformRuleSetFilter) []TransformRuleSet {
	out := make([]TransformRuleSet, 0, len(ruleSets))
	for _, item := range ruleSets {
		if f.KeyContains != "" && !strings.Contains(strings.ToLower(item.Key), strings.ToLower(f.KeyContains)) {
			continue
		}
		if f.Status != "" && item.Status != f.Status {
			continue
		}
		if f.Type != "" && item.Type != f.Type {
			continue
		}
		if f.ProcessingPurposeID != "" && item.ProcessingPurposeID != f.ProcessingPurposeID {
			continue
		}
		out = append(out, item)
	}
	sort.SliceStable(out, func(i, j int) bool {
		var less bool
		switch f.SortBy {
		case "updated_at":
			less = out[i].UpdatedAt.Before(out[j].UpdatedAt)
		case "created_at":
			less = out[i].CreatedAt.Before(out[j].CreatedAt)
		default:
			less = out[i].Key < out[j].Key
		}
		if f.Descending {
			return !less
		}
		return less
	})
	return out
}
func UniqueKeys(ruleSets []TransformRuleSet) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, f := range ruleSets {
		if !seen[f.Key] {
			seen[f.Key] = true
			out = append(out, f.Key)
		}
	}
	sort.Strings(out)
	return out
}
