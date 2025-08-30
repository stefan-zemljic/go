package z

import (
	"testing"
	"unsafe"
)

func TestWrapUnwrap(t *testing.T) {
	cases := [][]byte{
		[]byte("hello"),
		[]byte(""),
		[]byte("abcdef"),
	}
	for _, in := range cases {
		s := Wrap(in)
		out := Unwrap(s)
		if string(in) != s {
			t.Errorf("Wrap: got %q, want %q", s, in)
		}
		if string(out) != s {
			t.Errorf("Unwrap: got %q, want %q", out, s)
		}
	}
}

func TestAliasing(t *testing.T) {
	b := []byte("foo")
	s := Wrap(b)
	b[0] = 'b'
	if s != "boo" {
		t.Errorf("aliasing failed: got %q", s)
	}
}

func TestWrapUnwrapReferential(t *testing.T) {
	b := []byte("hello")
	s := Wrap(b)
	b2 := Unwrap(s)
	pb := unsafe.Pointer(&b[0])
	ps := unsafe.Pointer(unsafe.StringData(s))
	pb2 := unsafe.Pointer(&b2[0])
	if pb != ps || pb != pb2 {
		t.Errorf("expected shared backing, got different addresses: %p %p %p", pb, ps, pb2)
	}
}
