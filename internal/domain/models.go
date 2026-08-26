// Package implementation for privacy transformation and sensitive-value protection.
package domain

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrNotFound = errors.New("resource not found")
	ErrConflict = errors.New("resource conflict")
	ErrInvalid  = errors.New("invalid resource")
)

type ValueType string

const (
	TypeBool   ValueType = "boolean"
	TypeString ValueType = "string"
	TypeInt    ValueType = "integer"
	TypeJSON   ValueType = "json"
)

type PolicyWorkspace struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
type ProcessingPurpose struct {
	ID                string    `json:"id"`
	PolicyWorkspaceID string    `json:"workspace_id"`
	Name              string    `json:"name"`
	CreatedAt         time.Time `json:"created_at"`
}

type Rule struct {
	ID         string            `json:"id"`
	Priority   int               `json:"priority"`
	Tags       map[string]string `json:"tags,omitempty"`
	Percentage *int              `json:"percentage,omitempty"`
	StartAt    *time.Time        `json:"start_at,omitempty"`
	EndAt      *time.Time        `json:"end_at,omitempty"`
	Value      any               `json:"value"`
}

type TransformRuleSet struct {
	ID                      string    `json:"id"`
	PolicyWorkspaceID       string    `json:"workspace_id"`
	ProcessingPurposeID     string    `json:"purpose_id"`
	Key                     string    `json:"key"`
	Description             string    `json:"description,omitempty"`
	Type                    ValueType `json:"type"`
	DefaultValue            any       `json:"default_value"`
	Rules                   []Rule    `json:"rules,omitempty"`
	ActiveTransformRevision int       `json:"active_revision"`
	Status                  string    `json:"status"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}

type TransformRevision struct {
	Number             int    `json:"revision"`
	TransformRuleSetID string `json:"ruleSet_id"`
	Value              any    `json:"value"`
	Rules              []Rule `json:"rules,omitempty"`
	Status             string `json:"status"`
	CreatedAt          time.Time
	PublishedAt        *time.Time `json:"published_at,omitempty"`
}
type PolicyPublication struct {
	ID, TransformRuleSetID string
	TransformRevision      int
	ProcessingPurposeID    string
	Status                 string
	CreatedAt, UpdatedAt   time.Time
	Reason                 string
}

func ValidateValue(t ValueType, value any) error {
	switch t {
	case TypeBool:
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("%v: expected boolean", ErrInvalid)
		}
	case TypeString:
		if _, ok := value.(string); !ok {
			return fmt.Errorf("%v: expected string", ErrInvalid)
		}
	case TypeInt:
		switch value.(type) {
		case int, int64, float64, float32:
		default:
			return fmt.Errorf("%v: expected integer", ErrInvalid)
		}
	case TypeJSON:
		if value == nil {
			return fmt.Errorf("%v: JSON value cannot be nil", ErrInvalid)
		}
	default:
		return fmt.Errorf("%v: unsupported value type %q", ErrInvalid, t)
	}
	return nil
}

func (f *TransformRuleSet) Validate() error {
	if f.ID == "" || f.PolicyWorkspaceID == "" || f.ProcessingPurposeID == "" || f.Key == "" {
		return fmt.Errorf("%w: ruleSet identity is required", ErrInvalid)
	}
	if err := ValidateValue(f.Type, f.DefaultValue); err != nil {
		return err
	}
	return nil
}

func (v *TransformRevision) Validate(t ValueType) error {
	if v.Number < 1 {
		return fmt.Errorf("%w: revision must be positive", ErrInvalid)
	}
	return ValidateValue(t, v.Value)
}
