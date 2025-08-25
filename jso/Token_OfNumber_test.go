package jso

import (
	"testing"
)

func TestOfNumber_Valid(t *testing.T) {
	validCases := []string{
		"0", "123", "-456",
		"3.14", "-0.1",
		"1e10", "2E+5", "-3.4e-2",
	}
	for _, s := range validCases {
		t.Run(s, func(t *testing.T) {
			tok := OfNumber(s)
			if got := tok.Kind(); got != Number {
				t.Errorf("expected Number, got %v", got)
			}
			if tok.data.(number) != number(s) {
				t.Errorf("expected token data %q, got %q", s, tok.data)
			}
		})
	}
}
func TestOfNumber_Invalid(t *testing.T) {
	invalidCases := []string{
		"", "01", ".", "-", "e10", "abc",
	}
	for _, s := range invalidCases {
		t.Run(s, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected panic for invalid number %q", s)
				}
			}()
			_ = OfNumber(s)
		})
	}
}
