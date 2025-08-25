package jso

import (
	"testing"
)

func TestNull(t *testing.T) {
	(Token{nil}).Null()
	expectPanic(t, "Null", func() { (Token{true}).Null() })
}
func TestBool(t *testing.T) {
	if v := (Token{true}).Bool(); v != true {
		t.Errorf("expected true, got %v", v)
	}
	expectPanic(t, "Bool", func() { (Token{"nope"}).Bool() })
}
func TestInts(t *testing.T) {
	tok := Token{number("123")}
	if v := tok.Int(64); v != 123 {
		t.Errorf("Int failed, got %v", v)
	}
	if v := tok.Int8(8); v != 123 {
		t.Errorf("Int8 failed, got %v", v)
	}
	if v := tok.Int16(16); v != 123 {
		t.Errorf("Int16 failed, got %v", v)
	}
	if v := tok.Int32(32); v != 123 {
		t.Errorf("Int32 failed, got %v", v)
	}
	if v := tok.Int64(64); v != 123 {
		t.Errorf("Int64 failed, got %v", v)
	}
	expectPanic(t, "Int64 wrong type", func() { (Token{true}).Int64(64) })
	expectPanic(t, "Int64 bad parse", func() { (Token{number("abc")}).Int64(64) })
}
func TestUints(t *testing.T) {
	tok := Token{number("42")}
	if v := tok.Uint(64); v != 42 {
		t.Errorf("Uint failed, got %v", v)
	}
	if v := tok.Uint8(8); v != 42 {
		t.Errorf("Uint8 failed, got %v", v)
	}
	if v := tok.Uint16(16); v != 42 {
		t.Errorf("Uint16 failed, got %v", v)
	}
	if v := tok.Uint32(32); v != 42 {
		t.Errorf("Uint32 failed, got %v", v)
	}
	if v := tok.Uintptr(64); v != uintptr(42) {
		t.Errorf("Uintptr failed, got %v", v)
	}
	if v := tok.Uint64(64); v != 42 {
		t.Errorf("Uint64 failed, got %v", v)
	}
	expectPanic(t, "Uint64 wrong type", func() { (Token{"42"}).Uint64(64) })
	expectPanic(t, "Uint64 bad parse", func() { (Token{number("notnum")}).Uint64(64) })
}
func TestFloats(t *testing.T) {
	tok := Token{number("3.14")}
	if v := tok.Float32(32); v != float32(3.14) {
		t.Errorf("Float32 failed, got %v", v)
	}
	if v := tok.Float64(64); v != 3.14 {
		t.Errorf("Float64 failed, got %v", v)
	}
	expectPanic(t, "Float64 wrong type", func() { (Token{true}).Float64(64) })
	expectPanic(t, "Float64 bad parse", func() { (Token{number("oops")}).Float64(64) })
}
func TestNumber(t *testing.T) {
	if s := (Token{number("77")}).Number(); s != "77" {
		t.Errorf("expected 77, got %q", s)
	}
	expectPanic(t, "Number wrong type", func() { (Token{true}).Number() })
}
func TestString(t *testing.T) {
	if s := (Token{"hi"}).String(); s != "hi" {
		t.Errorf("expected hi, got %q", s)
	}
	expectPanic(t, "String wrong type", func() { (Token{123}).String() })
}
func TestArray(t *testing.T) {
	(Token{Array}).Array()
	expectPanic(t, "Array wrong type", func() { (Token{true}).Array() })
}
func TestObject(t *testing.T) {
	(Token{Object}).Object()
	expectPanic(t, "Object wrong type", func() { (Token{"oops"}).Object() })
}
func TestEnd(t *testing.T) {
	(Token{End}).End()
	expectPanic(t, "End wrong type", func() { (Token{123}).End() })
}
