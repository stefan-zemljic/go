package jso

import (
	"testing"
)

func TestRoundtrip(t *testing.T) {
	tests := []struct {
		name  string
		write func(*Writer)
		read  func(*Reader) string
		want  string
	}{
		{
			name:  "null",
			write: func(w *Writer) { w.Null() },
			read: func(r *Reader) string {
				if !r.Null() {
					t.Fatal("expected null")
				}
				return "null"
			},
			want: "null",
		},
		{
			name:  "bool true",
			write: func(w *Writer) { w.Bool(true) },
			read: func(r *Reader) string {
				if !r.MustBool() {
					t.Fatal("expected true")
				}
				return "true"
			},
			want: "true",
		},
		{
			name:  "int",
			write: func(w *Writer) { w.Int(42) },
			read: func(r *Reader) string {
				return string(r.MustNumber())
			},
			want: "42",
		},
		{
			name:  "float",
			write: func(w *Writer) { w.Float64(3.14) },
			read: func(r *Reader) string {
				return r.MustNumber()
			},
			want: "3.14",
		},
		{
			name:  "string",
			write: func(w *Writer) { w.String("abc") },
			read: func(r *Reader) string {
				return r.MustString()
			},
			want: "abc",
		},
		{
			name: "array",
			write: func(w *Writer) {
				w.Array()
				w.Int(1)
				w.Int(2)
				w.End()
			},
			read: func(r *Reader) string {
				if !r.Array() {
					t.Fatal("expected array")
				}
				a := r.MustInt()
				b := r.MustInt()
				r.MustEnd()
				return string([]rune{rune('0' + a), ',', rune('0' + b)})
			},
			want: "1,2",
		},
		{
			name: "object empty",
			write: func(w *Writer) {
				w.Object()
				w.End()
			},
			read: func(r *Reader) string {
				if !r.Object() {
					t.Fatal("expected object")
				}
				if !r.End() {
					t.Fatal("expected end")
				}
				return "{}"
			},
			want: "{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var w Writer
			tt.write(&w)

			r := &Reader{buf: []byte(w.Compact())}
			got := tt.read(r)
			if got != tt.want {
				t.Errorf("%s: got %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}
