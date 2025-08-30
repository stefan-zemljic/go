package z

import (
	"encoding/json"
	"testing"
)

func roundtrip(t *testing.T, input string) string {
	t.Helper()
	val, err := ParseJson([]byte(input))
	if err != nil {
		t.Fatalf("ParseJson(%q) failed: %v", input, err)
	}
	out, err := json.Marshal(val)
	if err != nil {
		t.Fatalf("Marshal(%q) failed: %v", input, err)
	}
	return string(out)
}

func TestParseJson_RoundtripStrings(t *testing.T) {
	cases := []string{
		"null",
		"true",
		"false",
		`"hello"`,
		"123",
		"-45.67",
		"[]",
		"[1,\"x\",true]",
		"{}",
		`{"a":1,"b":"x"}`,
	}
	for _, in := range cases {
		got := roundtrip(t, in)
		if got != in {
			t.Errorf("roundtrip failed: in=%q got=%q", in, got)
		}
	}
}

func TestParseJson_Errors(t *testing.T) {
	cases := []string{
		"{bad}",
		"{",
		"[",
		`{"a":1`,
		`[1,2`,
		`{1:2}`,
		`[falsy]`,
	}
	for _, in := range cases {
		if _, err := ParseJson([]byte(in)); err == nil {
			t.Errorf("expected error for input %q, got nil", in)
		}
	}
}
