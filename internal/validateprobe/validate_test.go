package validateprobe

import (
	"errors"
	"testing"

	"github.com/ali/go-0821/privacy-transform-service/internal/domain"
)

func TestValidateBoolInvalidChain(t *testing.T) {
	if err := domain.ValidateValue(domain.TypeBool, 123); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestValidateStringInvalidChain(t *testing.T) {
	if err := domain.ValidateValue(domain.TypeString, 123); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestValidateIntegerInvalidChain(t *testing.T) {
	if err := domain.ValidateValue(domain.TypeInt, "x"); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestValidateJSONNilChain(t *testing.T) {
	if err := domain.ValidateValue(domain.TypeJSON, nil); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestValidateUnsupportedTypeChain(t *testing.T) {
	if err := domain.ValidateValue(domain.ValueType("bogus"), "x"); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestValidateRuleInvalidPathChain(t *testing.T) {
	if err := domain.ValidateRule(domain.FieldRule{Path: "", Action: "mask"}); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}

func TestValidateRuleUnsupportedActionChain(t *testing.T) {
	if err := domain.ValidateRule(domain.FieldRule{Path: "email", Action: "scramble"}); !errors.Is(err, domain.ErrInvalid) {
		t.Fatalf("errors.Is(err, ErrInvalid) = false, err=%v", err)
	}
}
