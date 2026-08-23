// Package implementation for privacy transformation and sensitive-value protection.
package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Classification struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Level       int    `json:"level"`
	Description string `json:"description,omitempty"`
}
type FieldRule struct {
	Path           string `json:"path"`
	Classification string `json:"classification"`
	Action         string `json:"action"`
	Salt           string `json:"salt,omitempty"`
	Preserve       int    `json:"preserve,omitempty"`
}
type TransformPolicy struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	Version   int         `json:"version"`
	Status    string      `json:"status"`
	Rules     []FieldRule `json:"rules"`
	UpdatedAt time.Time   `json:"updated_at"`
}
type TransformResult struct {
	RequestID string         `json:"request_id"`
	PolicyID  string         `json:"policy_id"`
	Data      map[string]any `json:"data"`
	Summary   Summary        `json:"summary"`
	CreatedAt time.Time      `json:"created_at"`
}
type Summary struct {
	Processed   int `json:"processed"`
	Transformed int `json:"transformed"`
	Deleted     int `json:"deleted"`
	Hashes      int `json:"hashes"`
	Tokens      int `json:"tokens"`
}

var pathPart = regexp.MustCompile(`^[A-Za-z0-9_.-]+$`)

func ValidateRule(r FieldRule) error {
	if r.Path == "" || (!pathPart.MatchString(r.Path) && !strings.Contains(r.Path, "[")) {
		return fmt.Errorf("%w: invalid field path", ErrInvalid)
	}
	switch r.Action {
	case "mask", "hash", "tokenize", "truncate", "generalize", "delete":
	default:
		return fmt.Errorf("%w: unsupported action", ErrInvalid)
	}
	return nil
}
func HashValue(value string, salt string) string {
	h := sha256.Sum256([]byte(salt + value))
	return hex.EncodeToString(h[:])
}
func MaskValue(value string, preserve int) string {
	if preserve < 0 {
		preserve = 0
	}
	r := []rune(value)
	if preserve*2 >= len(r) {
		return strings.Repeat("*", len(r))
	}
	return string(r[:preserve]) + strings.Repeat("*", len(r)-preserve*2) + string(r[len(r)-preserve:])
}
