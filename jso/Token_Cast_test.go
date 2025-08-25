package jso

import (
	"testing"
)

func TestAsNull(t *testing.T) {
	if !(Token{nil}).AsNull() {
		t.Error("expected AsNull true for nil")
	}
	if (Token{42}).AsNull() {
		t.Error("expected AsNull false for non-nil")
	}
}
func TestAsBool(t *testing.T) {
	if v, ok := (Token{true}).AsBool(); !ok || v != true {
		t.Errorf("expected (true,true), got (%v,%v)", v, ok)
	}
	if v, ok := (Token{"nope"}).AsBool(); ok || v {
		t.Errorf("expected (false,false), got (%v,%v)", v, ok)
	}
}
func TestAsInts(t *testing.T) {
	tok := Token{number("123")}
	if v, ok := tok.AsInt(64); !ok || v != 123 {
		t.Errorf("AsInt failed, got (%v,%v)", v, ok)
	}
	if v, ok := tok.AsInt8(8); !ok || v != 123 {
		t.Errorf("AsInt8 failed, got (%v,%v)", v, ok)
	}
	if _, ok := (Token{number("bad")}).AsInt(64); ok {
		t.Error("expected parse failure for bad int")
	}
	if v, ok := (Token{"123"}).AsInt(64); ok || v != 0 {
		t.Error("expected (0,false) for wrong type")
	}
}
func TestAsUints(t *testing.T) {
	tok := Token{number("42")}
	if v, ok := tok.AsUint(64); !ok || v != 42 {
		t.Errorf("AsUint failed, got (%v,%v)", v, ok)
	}
	if v, ok := tok.AsUint8(8); !ok || v != 42 {
		t.Errorf("AsUint8 failed, got (%v,%v)", v, ok)
	}
	if _, ok := (Token{number("-5")}).AsUint(64); ok {
		t.Error("expected parse failure for negative")
	}
	if _, ok := (Token{"42"}).AsUint(64); ok {
		t.Error("expected false for wrong type")
	}
}
func TestAsOtherIntsAndUints(t *testing.T) {
	tok := Token{number("123")}
	if v, ok := tok.AsInt16(16); !ok || v != 123 {
		t.Errorf("AsInt16 failed, got (%v,%v)", v, ok)
	}
	if v, ok := tok.AsInt32(32); !ok || v != 123 {
		t.Errorf("AsInt32 failed, got (%v,%v)", v, ok)
	}
	if v, ok := tok.AsInt64(64); !ok || v != 123 {
		t.Errorf("AsInt64 failed, got (%v,%v)", v, ok)
	}
	if v, ok := tok.AsUint16(16); !ok || v != 123 {
		t.Errorf("AsUint16 failed, got (%v,%v)", v, ok)
	}
	if v, ok := tok.AsUint32(32); !ok || v != 123 {
		t.Errorf("AsUint32 failed, got (%v,%v)", v, ok)
	}
	if v, ok := tok.AsUint64(64); !ok || v != 123 {
		t.Errorf("AsUint64 failed, got (%v,%v)", v, ok)
	}
	if v, ok := tok.AsUintptr(64); !ok || v != uintptr(123) {
		t.Errorf("AsUintptr failed, got (%v,%v)", v, ok)
	}
}
func TestAsFloats(t *testing.T) {
	tok := Token{number("3.14")}
	if v, ok := tok.AsFloat32(32); !ok || v != 3.14 {
		t.Errorf("AsFloat32 failed, got (%v,%v)", v, ok)
	}
	if v, ok := tok.AsFloat64(64); !ok || v != 3.14 {
		t.Errorf("AsFloat64 failed, got (%v,%v)", v, ok)
	}
	if _, ok := (Token{number("NaNstr")}).AsFloat64(64); ok {
		t.Error("expected parse failure for bad float")
	}
	if _, ok := (Token{true}).AsFloat32(32); ok {
		t.Error("expected false for wrong type")
	}
}
func TestAsNumber(t *testing.T) {
	if s, ok := (Token{number("55")}).AsNumber(); !ok || s != "55" {
		t.Errorf("expected (55,true), got (%q,%v)", s, ok)
	}
	if s, ok := (Token{"55"}).AsNumber(); ok || s != "" {
		t.Errorf("expected (\"\",false), got (%q,%v)", s, ok)
	}
}
func TestAsString(t *testing.T) {
	if s, ok := (Token{"hello"}).AsString(); !ok || s != "hello" {
		t.Errorf("expected (hello,true), got (%q,%v)", s, ok)
	}
	if s, ok := (Token{number("1")}).AsString(); ok || s != "" {
		t.Errorf("expected (\"\",false), got (%q,%v)", s, ok)
	}
}
func TestAsKinds(t *testing.T) {
	if !(Token{Array}).AsArray() {
		t.Error("expected AsArray true")
	}
	if (Token{"x"}).AsArray() {
		t.Error("expected AsArray false")
	}
	if !(Token{Object}).AsObject() {
		t.Error("expected AsObject true")
	}
	if (Token{123}).AsObject() {
		t.Error("expected AsObject false")
	}
	if !(Token{End}).AsEnd() {
		t.Error("expected AsEnd true")
	}
	if (Token{nil}).AsEnd() {
		t.Error("expected AsEnd false")
	}
}
