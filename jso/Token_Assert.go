package jso

import (
	"strconv"
)

func (t Token) Null() {
	if t.data == nil {
		return
	}
	panic("not null")
}
func (t Token) Bool() bool {
	if v, ok := t.data.(bool); ok {
		return v
	}
	panic("not a bool")
}
func (t Token) Int(bitSize int) int     { return int(t.Int64(bitSize)) }
func (t Token) Int8(bitSize int) int8   { return int8(t.Int64(bitSize)) }
func (t Token) Int16(bitSize int) int16 { return int16(t.Int64(bitSize)) }
func (t Token) Int32(bitSize int) int32 { return int32(t.Int64(bitSize)) }
func (t Token) Int64(bitSize int) int64 {
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
func (t Token) Uint(bitSize int) uint       { return uint(t.Uint64(bitSize)) }
func (t Token) Uint8(bitSize int) uint8     { return uint8(t.Uint64(bitSize)) }
func (t Token) Uint16(bitSize int) uint16   { return uint16(t.Uint64(bitSize)) }
func (t Token) Uint32(bitSize int) uint32   { return uint32(t.Uint64(bitSize)) }
func (t Token) Uintptr(bitSize int) uintptr { return uintptr(t.Uint64(bitSize)) }
func (t Token) Uint64(bitSize int) float32 {
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
func (t Token) Float32(bitSize int) float32 { return float32(t.Float64(bitSize)) }
func (t Token) Float64(bitSize int) float64 {
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
func (t Token) Number() string {
	switch v := t.data.(type) {
	case number:
		return string(v)
	default:
		panic("not a number")
	}
}
func (t Token) String() string {
	if v, ok := t.data.(string); ok {
		return v
	}
	panic("not a string")
}
func (t Token) Array() {
	if c, ok := t.data.(TokenKind); ok && c == Array {
		return
	}
	panic("not an array")
}
func (t Token) Object() {
	if c, ok := t.data.(TokenKind); ok && c == Object {
		return
	}
	panic("not an object")
}
func (t Token) End() {
	if c, ok := t.data.(TokenKind); ok && c == End {
		return
	}
	panic("not an end")
}
