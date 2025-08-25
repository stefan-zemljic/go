package jso

import (
	"testing"
)

func TestTokenOf_Primitives(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected TokenKind
	}{
		{"nil", nil, Null},
		{"bool", true, Bool},
		{"int", 42, Number},
		{"int8", int8(8), Number},
		{"int16", int16(16), Number},
		{"int32", int32(32), Number},
		{"int64", int64(64), Number},
		{"uint", uint(42), Number},
		{"uint8", uint8(8), Number},
		{"uint16", uint16(16), Number},
		{"uint32", uint32(32), Number},
		{"uint64", uint64(64), Number},
		{"float32", float32(3.14), Number},
		{"float64", 6.28, Number},
		{"string", "hello", String},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := TokenOf(tt.input)
			if got := tok.Kind(); got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
func TestTokenOf_TokenKindAllowed(t *testing.T) {
	for _, kind := range []TokenKind{Array, Object, End} {
		tok := TokenOf(kind)
		if got := tok.Kind(); got != kind {
			t.Errorf("expected %v, got %v", kind, got)
		}
	}
}
func TestTokenOf_TokenKindDisallowed(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for disallowed TokenKind")
		}
	}()
	_ = TokenOf(Null)
}
func TestTokenOf_UnsupportedType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for unsupported type")
		}
	}()
	_ = TokenOf(struct{}{})
}
func TestTokenKind_InvalidTokenKind(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for invalid token kind")
		}
	}()
	tok := Token{data: String}
	_ = tok.Kind()
}
func TestTokenKind_InvalidType(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for invalid type")
		}
	}()
	tok := Token{data: []byte("oops")}
	_ = tok.Kind()
}
