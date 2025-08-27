package jso

import (
	"strconv"
)

func (t Token) MustNull() {
	if t.data == nil {
		return
	}
	panic("not null")
}
func (t Token) MustBool() bool {
	if v, ok := t.data.(bool); ok {
		return v
	}
	panic("not a bool")
}
func (t Token) MustInt(bitSize int) int     { return int(t.MustInt64(bitSize)) }
func (t Token) MustInt8(bitSize int) int8   { return int8(t.MustInt64(bitSize)) }
func (t Token) MustInt16(bitSize int) int16 { return int16(t.MustInt64(bitSize)) }
func (t Token) MustInt32(bitSize int) int32 { return int32(t.MustInt64(bitSize)) }
func (t Token) MustInt64(bitSize int) int64 {
	switch v := t.data.(type) {
	case number:
		i, err := strconv.ParseInt(string(v), 10, bitSize)
		if err != nil {
			panic(err)
		}
		return i
	default:
		panic("not a number")
	}
}
func (t Token) MustUint(bitSize int) uint       { return uint(t.MustUint64(bitSize)) }
func (t Token) MustUint8(bitSize int) uint8     { return uint8(t.MustUint64(bitSize)) }
func (t Token) MustUint16(bitSize int) uint16   { return uint16(t.MustUint64(bitSize)) }
func (t Token) MustUint32(bitSize int) uint32   { return uint32(t.MustUint64(bitSize)) }
func (t Token) MustUintptr(bitSize int) uintptr { return uintptr(t.MustUint64(bitSize)) }
func (t Token) MustUint64(bitSize int) float32 {
	switch v := t.data.(type) {
	case number:
		u, err := strconv.ParseUint(string(v), 10, bitSize)
		if err != nil {
			panic(err)
		}
		return float32(u)
	default:
		panic("not a number")
	}
}
func (t Token) MustFloat32(bitSize int) float32 { return float32(t.MustFloat64(bitSize)) }
func (t Token) MustFloat64(bitSize int) float64 {
	switch v := t.data.(type) {
	case number:
		f, err := strconv.ParseFloat(string(v), bitSize)
		if err != nil {
			panic(err)
		}
		return f
	default:
		panic("not a number")
	}
}
func (t Token) MustNumber() string {
	switch v := t.data.(type) {
	case number:
		return string(v)
	default:
		panic("not a number")
	}
}
func (t Token) MustString() string {
	if v, ok := t.data.(string); ok {
		return v
	}
	panic("not a string")
}
func (t Token) MustArr() {
	if c, ok := t.data.(TokenKind); ok && c == Arr {
		return
	}
	panic("not an array")
}
func (t Token) MustObj() {
	if c, ok := t.data.(TokenKind); ok && c == Obj {
		return
	}
	panic("not an object")
}
func (t Token) MustEnd() {
	if c, ok := t.data.(TokenKind); ok && c == End {
		return
	}
	panic("not an end")
}
