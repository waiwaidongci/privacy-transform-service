// Package implementation for privacy transformation and sensitive-value protection.
package repository

import (
	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
	"strings"
)

type TransformRuleSetIndex struct {
	byPolicyWorkspace map[string][]string
	byKey             map[string]string
}

func NewTransformRuleSetIndex() *TransformRuleSetIndex {
	return &TransformRuleSetIndex{byPolicyWorkspace: map[string][]string{}, byKey: map[string]string{}}
}
func (i *TransformRuleSetIndex) Add(f domain.TransformRuleSet) {
	i.byPolicyWorkspace[f.PolicyWorkspaceID] = append(i.byPolicyWorkspace[f.PolicyWorkspaceID], f.ID)
	i.byKey[strings.ToLower(f.Key)] = f.ID
}
func (i *TransformRuleSetIndex) Remove(f domain.TransformRuleSet) {
	ids := i.byPolicyWorkspace[f.PolicyWorkspaceID]
	out := ids[:0]
	for _, id := range ids {
		if id != f.ID {
			out = append(out, id)
		}
	}
	i.byPolicyWorkspace[f.PolicyWorkspaceID] = out
	delete(i.byKey, f.Key)
}
func (i *TransformRuleSetIndex) ListByWorkspace(workspaceID string) []string {
	return i.byPolicyWorkspace[workspaceID]
}
func (i *TransformRuleSetIndex) FindByKey(key string) string { return i.byKey[strings.ToLower(key)] }
