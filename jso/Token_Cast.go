package jso

import (
	"strconv"
)

func (t Token) AsNull() bool {
	if t.data == nil {
		return true
	}
	return false
}
func (t Token) AsBool() (bool, bool) {
	if v, ok := t.data.(bool); ok {
		return v, true
	}
	return false, false
}
func (t Token) AsInt(bitSize int) (int, bool)     { return asInt[int](t, bitSize) }
func (t Token) AsInt8(bitSize int) (int8, bool)   { return asInt[int8](t, bitSize) }
func (t Token) AsInt16(bitSize int) (int16, bool) { return asInt[int16](t, bitSize) }
func (t Token) AsInt32(bitSize int) (int32, bool) { return asInt[int32](t, bitSize) }
func (t Token) AsInt64(bitSize int) (int64, bool) { return asInt[int64](t, bitSize) }
func asInt[T interface {
	int | int8 | int16 | int32 | int64
}](t Token, bitSize int) (T, bool) {
	switch v := t.data.(type) {
	case number:
		i, err := strconv.ParseInt(string(v), 10, bitSize)
		if err != nil {
			return 0, false
		}
		return T(i), true
	default:
		return 0, false
	}
}
func (t Token) AsUint(bitSize int) (uint, bool)       { return genUint[uint](t, bitSize) }
func (t Token) AsUint8(bitSize int) (uint8, bool)     { return genUint[uint8](t, bitSize) }
func (t Token) AsUint16(bitSize int) (uint16, bool)   { return genUint[uint16](t, bitSize) }
func (t Token) AsUint32(bitSize int) (uint32, bool)   { return genUint[uint32](t, bitSize) }
func (t Token) AsUint64(bitSize int) (uint64, bool)   { return genUint[uint64](t, bitSize) }
func (t Token) AsUintptr(bitSize int) (uintptr, bool) { return genUint[uintptr](t, bitSize) }
func genUint[T interface {
	uint | uint8 | uint16 | uint32 | uint64 | uintptr
}](t Token, bitSize int) (T, bool) {
	switch v := t.data.(type) {
	case number:
		u, err := strconv.ParseUint(string(v), 10, bitSize)
		if err != nil {
			return 0, false
		}
		return T(u), true
	default:
		return 0, false
	}
}
func (t Token) AsFloat32(bitSize int) (float32, bool) { return genFloat[float32](t, bitSize) }
func (t Token) AsFloat64(bitSize int) (float64, bool) { return genFloat[float64](t, bitSize) }
func genFloat[T interface{ float32 | float64 }](t Token, bitSize int) (T, bool) {
	switch v := t.data.(type) {
	case number:
		f, err := strconv.ParseFloat(string(v), bitSize)
		if err != nil {
			return 0, false
		}
		return T(f), true
	default:
		return 0, false
	}
}
func (t Token) AsNumber() (string, bool) {
	switch v := t.data.(type) {
	case number:
		return string(v), true
	default:
		return "", false
	}
}
func (t Token) AsString() (string, bool) {
	if v, ok := t.data.(string); ok {
		return v, true
	}
	return "", false
}
func (t Token) AsArray() bool {
	if c, ok := t.data.(TokenKind); ok && c == Array {
		return true
	}
	return false
}
func (t Token) AsObject() bool {
	if c, ok := t.data.(TokenKind); ok && c == Object {
		return true
	}
	return false
}
func (t Token) AsEnd() bool {
	if c, ok := t.data.(TokenKind); ok && c == End {
		return true
	}
	return false
}
