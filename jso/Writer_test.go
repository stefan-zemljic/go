package jso

import (
	"errors"
	"fmt"
	"io"
	"testing"
)

func TestWriter_Success(t *testing.T) {
	type test struct {
		name string
		f    func(w *Writer)
		want string
	}
	tests := []test{
		{"null", func(w *Writer) { w.Null() }, "null"},
		{"bool", func(w *Writer) { w.Bool(true) }, "true"},
		{"int", func(w *Writer) { w.Int(-42) }, "-42"},
		{"int8", func(w *Writer) { w.Int8(8) }, "8"},
		{"int16", func(w *Writer) { w.Int16(16) }, "16"},
		{"int32", func(w *Writer) { w.Int32(32) }, "32"},
		{"int64", func(w *Writer) { w.Int64(64) }, "64"},
		{"uint", func(w *Writer) { w.Uint(7) }, "7"},
		{"uint8", func(w *Writer) { w.Uint8(8) }, "8"},
		{"uint16", func(w *Writer) { w.Uint16(16) }, "16"},
		{"uint32", func(w *Writer) { w.Uint32(32) }, "32"},
		{"uint64", func(w *Writer) { w.Uint64(64) }, "64"},
		{"uintptr", func(w *Writer) { w.Uintptr(9) }, "9"},
		{"float32", func(w *Writer) { w.Float32(1.5) }, "1.5"},
		{"float64", func(w *Writer) { w.Float64(2.5) }, "2.5"},
		{"number", func(w *Writer) { w.Number("123.45") }, "123.45"},
		{"string", func(w *Writer) { w.String("abc") }, `"abc"`},
		{"array", func(w *Writer) {
			w.Array()
			w.Int(1)
			w.Int(2)
			w.End()
		}, "[1,2]"},
		{"object", func(w *Writer) {
			w.Object()
			w.String("a")
			w.Int(1)
			w.String("b")
			w.Bool(false)
			w.End()
		}, `{"a":1,"b":false}`},
		{"nested", func(w *Writer) {
			w.Object()
			w.String("arr")
			w.Array()
			w.String("x")
			w.String("y")
			w.End()
			w.End()
		}, `{"arr":["x","y"]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w Writer
			tt.f(&w)
			if got := w.Compact(); got != tt.want {
				t.Errorf("Writer.%s() = %v, want %v", tt.name, got, tt.want)
			}
			w.CompactTo(io.Discard)
			w.CompactBytes()
			w.Pretty()
			w.PrettyTo(io.Discard)
			w.Indent("", "  ")
			w.IndentTo(io.Discard, "", "  ")

		})
	}
}

func mustPanic(t *testing.T, f func(), substr string) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("expected panic %q, but no panic", substr)
		} else if msg := fmt.Sprintf("%v", r); msg != substr {
			t.Fatalf("expected panic %q, got %q", substr, msg)
		}
	}()
	f()
}

type badWriter struct{}

func (badWriter) Write(p []byte) (int, error) { return 0, errors.New("write failed") }

func TestWriter_Errors(t *testing.T) {
	tests := []struct {
		name    string
		f       func()
		wantErr string
	}{
		{
			name: "End with no array/object",
			f: func() {
				var w Writer
				w.End()
			},
			wantErr: "jso: cannot end, no open array or object",
		},
		{
			name: "End at top-level complete (atEnd)",
			f: func() {
				var w Writer
				w.Null()
				w.End()
			},
			wantErr: "jso: cannot end, no open array or object",
		},
		{
			name: "End atStart (never wrote anything)",
			f: func() {
				var w Writer
				// simulate atStart → push into states then reset manually
				w.states = []state{atStart}
				w.state = atStart
				w.End()
			},
			wantErr: "jso: cannot end, not in array or object",
		},
		{
			name: "End afterObjectKey",
			f: func() {
				var w Writer
				w.Object()
				w.String("a") // moves state → afterObjectKey
				w.End()
			},
			wantErr: "jso: cannot end after object key",
		},
		{
			name: "Object key must be string (atObjectStart)",
			f: func() {
				var w Writer
				w.Object()
				w.Int(1)
			},
			wantErr: "jso: object keys must be strings",
		},
		{
			name: "Object key must be string (afterObjectValue)",
			f: func() {
				var w Writer
				w.Object()
				w.String("a")
				w.Int(1)
				w.Int(2)
			},
			wantErr: "jso: object keys must be strings",
		},
		{
			name: "Cannot write after document complete",
			f: func() {
				var w Writer
				w.Bool(true)
				w.Int(1)
			},
			wantErr: "jso: cannot write value, document already complete",
		},
		{
			name: "PrettyBytes invalid JSON",
			f: func() {
				var w Writer
				w.Object() // never closed
				w.PrettyBytes()
			},
			wantErr: "unexpected end of JSON input",
		},
		{
			name: "Indent invalid JSON",
			f: func() {
				var w Writer
				w.Array()
				w.String("x")
				// missing End
				w.Indent("", "  ")
			},
			wantErr: "unexpected end of JSON input",
		},
		{
			name: "CompactTo writer error",
			f: func() {
				var w Writer
				w.Bool(true)
				w.CompactTo(badWriter{})
			},
			wantErr: "write failed",
		},
		{
			name: "PrettyTo writer error",
			f: func() {
				var w Writer
				w.Bool(true)
				w.PrettyTo(badWriter{})
			},
			wantErr: "write failed",
		},
		{
			name: "IndentTo writer error",
			f: func() {
				var w Writer
				w.Bool(true)
				w.IndentTo(badWriter{}, "", "  ")
			},
			wantErr: "write failed",
		},
		{
			name: "IndentBytes invalid JSON",
			f: func() {
				var w Writer
				w.Object()
				w.IndentBytes("", "  ")
			},
			wantErr: "unexpected end of JSON input",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Fatalf("expected panic %q, got none", tt.wantErr)
				} else {
					got := fmt.Sprint(r)
					if got != tt.wantErr {
						t.Fatalf("expected panic %q, got %q", tt.wantErr, got)
					}
				}
			}()
			tt.f()
		})
	}
}
