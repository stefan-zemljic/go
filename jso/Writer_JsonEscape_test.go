package jso

import (
	"testing"
)

func TestWriter_JsonEscape(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", `""`},
		{"backslash", `\`, `"\\"`},
		{"quote", `"`, `"\""`},
		{"backspace", "\b", `"\b"`},
		{"tab", "\t", `"\t"`},
		{"newline", "\n", `"\n"`},
		{"formfeed", "\f", `"\f"`},
		{"carriage return", "\r", `"\r"`},
		{"control char <0x20", "\x01", `"\u0001"`},
		{"non-BMP rune (>0xFFFF)", string(rune(0x1F600)), `"\ud83d\ude00"`},
		{"printable rune", "A", `"A"`},
		{"mixed string", "A\"\n", `"A\"\n"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w Writer
			w.JsonEscape(tt.input)
			got := string(w.Buffer)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
