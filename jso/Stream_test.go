package jso

import (
	"testing"
)

func TestStream_AddAndAddAll(t *testing.T) {
	var s Stream
	s.Add(nil)
	if len(s.Tokens) != 1 {
		t.Errorf("expected 1 token, got %d", len(s.Tokens))
	}
	s.AddAll(true, 42, "hi")
	if len(s.Tokens) != 4 {
		t.Errorf("expected 4 tokens total, got %d", len(s.Tokens))
	}
	if v, ok := s.Tokens[1].AsBool(); !ok || v != true {
		t.Errorf("expected bool true, got %v %v", v, ok)
	}
	if v, ok := s.Tokens[2].AsInt(0); !ok || v != 42 {
		t.Errorf("expected int 42, got %v %v", v, ok)
	}
	if v, ok := s.Tokens[3].AsString(); !ok || v != "hi" {
		t.Errorf("expected string 'hi', got %v %v", v, ok)
	}
}
func TestStream_String(t *testing.T) {
	tests := []struct {
		name   string
		values []any
		want   string
	}{
		{"nil", []any{nil}, "null"},
		{"bool true", []any{true}, "true"},
		{"bool false", []any{false}, "false"},
		{"number", []any{123}, "123"},
		{"string", []any{"hi"}, `"hi"`},
		{"array", []any{Array, End}, "[]"},
		{"object", []any{Object, End}, "{}"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Stream
			s.AddAll(tt.values...)
			got := s.String()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
func TestStream_String_InvalidTokenKind(t *testing.T) {
	var s Stream
	s.Tokens = []Token{{TokenKind(-1)}}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for invalid token kind")
		}
	}()
	_ = s.String()
}
func TestStream_PrettyAndIndent(t *testing.T) {
	var s Stream
	s.Tokens = []Token{TokenOf(Array), TokenOf(1), TokenOf(true), TokenOf(End)}
	pretty := s.Pretty()
	if pretty != "[\n  1,\n  true\n]" {
		t.Errorf("unexpected pretty output: %s", pretty)
	}
	indented := s.Indent("  ", "  ")
	if indented != "[\n    1,\n    true\n  ]" {
		t.Errorf("unexpected indented output: %s", indented)
	}
}
func TestStream_PeekBranches(t *testing.T) {
	var s Stream
	if tok, ok := s.Peek(); ok || (tok != Token{}) {
		t.Fatalf("expected empty Peek, got %v %v", tok, ok)
	}
	s.Tokens = []Token{TokenOf("x")}
	if tok, ok := s.Peek(); !ok || tok == (Token{}) {
		t.Fatalf("expected non-empty Peek, got %v %v", tok, ok)
	}
}
func TestStream_ReadBranches(t *testing.T) {
	var s Stream
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on Read with empty stream")
		}
	}()
	_ = s.Read()
}
func TestStream_ReadNormal(t *testing.T) {
	s := Stream{Tokens: []Token{TokenOf("a"), TokenOf("b")}}
	tok1 := s.Read()
	if tok1 == (Token{}) || len(s.Tokens) != 1 {
		t.Fatalf("expected to read first token, got %v, remaining %d", tok1, len(s.Tokens))
	}
	tok2 := s.Read()
	if tok2 == (Token{}) || !s.IsEmpty() {
		t.Fatalf("expected to read second token, got %v, IsEmpty=%v", tok2, s.IsEmpty())
	}
}
func TestStream_IsEmptyAndMore(t *testing.T) {
	var s Stream
	if !s.IsEmpty() || s.More() {
		t.Fatalf("expected empty stream")
	}
	s.Tokens = []Token{TokenOf(1)}
	if s.IsEmpty() || !s.More() {
		t.Fatalf("expected non-empty stream")
	}
}
