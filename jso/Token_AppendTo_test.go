package jso

import (
	"testing"
)

func TestAppendTo_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		token       Token
		setupBuf    func(*Writer)
		want        string
		shouldPanic bool
	}{
		{
			name:     "End inside object",
			token:    Token{End},
			setupBuf: func(b *Writer) { b.PushState(AtObjectStart) },
			want:     "}",
		},
		{
			name:     "End inside array",
			token:    Token{End},
			setupBuf: func(b *Writer) { b.PushState(AtArrayStart) },
			want:     "]",
		},
		{
			name:        "End with no state",
			token:       Token{End},
			shouldPanic: true,
		},
		{
			name:  "nil",
			token: Token{nil},
			want:  "null",
		},
		{
			name:  "bool true",
			token: Token{true},
			want:  "true",
		},
		{
			name:  "bool false",
			token: Token{false},
			want:  "false",
		},
		{
			name:  "number",
			token: Token{number("123")},
			want:  "123",
		},
		{
			name:  "string",
			token: Token{"hi"},
			want:  `"hi"`,
		},
		{
			name:  "array token",
			token: Token{Array},
			want:  "[",
		},
		{
			name:  "object token",
			token: Token{Object},
			want:  "{",
		},
		{
			name:        "invalid TokenKind",
			token:       Token{TokenKind(-1)},
			shouldPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w Writer
			if tt.setupBuf != nil {
				tt.setupBuf(&w)
			}
			if tt.shouldPanic {
				expectPanic(t, tt.name, func() {
					tt.token.AddTo(&w)
				})
				return
			}
			tt.token.AddTo(&w)
			if got := string(w.Buffer); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
