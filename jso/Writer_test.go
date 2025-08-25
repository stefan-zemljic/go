package jso

import (
	"io"
	"strings"
	"testing"
)

func TestWriter_AppendAndAppendString(t *testing.T) {
	var b Writer
	b.Append('a', 'b')
	if got := b.String(); got != "ab" {
		t.Errorf("expected %q, got %q", "ab", got)
	}
	b.AppendString("cd")
	if got := b.String(); got != "abcd" {
		t.Errorf("expected %q, got %q", "abcd", got)
	}
}
func TestWriter_Write(t *testing.T) {
	var b Writer
	data := []byte("hello")
	n, err := b.Write(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != len(data) {
		t.Errorf("expected written len %d, got %d", len(data), n)
	}
	if got := b.String(); got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
	var _ io.Writer = &b
}
func TestWriter_Append(t *testing.T) {
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
			var b Writer
			b.AddAll(tt.values...)
			if got := b.String(); got != tt.want {
				t.Errorf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
func TestWriter_Indent_SuccessAndPanic(t *testing.T) {
	w := &Writer{}
	w.Buffer = []byte(`{"a":1}`)
	got := w.Indent("", "  ")
	if !strings.Contains(got, "\n") {
		t.Fatalf("expected pretty JSON with newline, got %q", got)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic on invalid JSON")
		}
	}()
	bad := &Writer{}
	bad.Buffer = []byte(`{invalid}`)
	_ = bad.Indent("", "  ")
}
func TestWriter_Reset(t *testing.T) {
	w := &Writer{
		Buffer:     []byte("abc"),
		State:      123,
		PrevStates: []WriterState{1, 2, 3},
	}
	w.Reset()
	if len(w.Buffer) != 0 {
		t.Fatalf("expected Buffer cleared, got %v", w.Buffer)
	}
	if w.State != AtStart {
		t.Fatalf("expected State=AtStart, got %v", w.State)
	}
	if len(w.PrevStates) != 0 {
		t.Fatalf("expected PrevStates cleared, got %v", w.PrevStates)
	}
}
func TestStringOf_Pretty_Indent(t *testing.T) {
	s1 := StringOf(Array, 123, "abc", End)
	if s1 != `[123,"abc"]` {
		t.Fatalf("StringOf missing expected content: %q", s1)
	}
	s2 := Pretty(Object, "x", 1, "y", 2, End)
	if !strings.Contains(s2, "\n") {
		t.Fatalf("Pretty output should be indented JSON: %q", s2)
	}
	s3 := Indent("", "  ", Object, "x", 1, "y", 2, End)
	if !strings.Contains(s3, "\n") || !strings.Contains(s3, "y") {
		t.Fatalf("Indent output mismatch: %q", s3)
	}
}
