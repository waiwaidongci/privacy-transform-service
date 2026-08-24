package domain

import (
	"encoding/json"
	"fmt"
)

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
