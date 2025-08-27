package jso

func (s *Data) MustNull() {
	assertNextToken[struct{}](s, func(t Token) struct{} {
		t.MustNull()
		return struct{}{}
	})
}
func (s *Data) MustBool() bool              { return assertNextToken[bool](s, Token.MustBool) }
func (s *Data) MustInt(bitSize int) int     { return mustInt[int](s, bitSize) }
func (s *Data) MustInt8(bitSize int) int8   { return mustInt[int8](s, bitSize) }
func (s *Data) MustInt16(bitSize int) int16 { return mustInt[int16](s, bitSize) }
func (s *Data) MustInt32(bitSize int) int32 { return mustInt[int32](s, bitSize) }
func (s *Data) MustInt64(bitSize int) int64 { return mustInt[int64](s, bitSize) }
func mustInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](d *Data, bitSize int) T {
	return assertNextToken[T](d, tokenMustInt[T](bitSize))
}
func tokenMustInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](bitSize int) func(t Token) T {
	return func(t Token) T { return T(t.MustInt64(bitSize)) }
}
func (s *Data) MustUint(bitSize int) uint       { return mustUint[uint](s, bitSize) }
func (s *Data) MustUint8(bitSize int) uint8     { return mustUint[uint8](s, bitSize) }
func (s *Data) MustUint16(bitSize int) uint16   { return mustUint[uint16](s, bitSize) }
func (s *Data) MustUint32(bitSize int) uint32   { return mustUint[uint32](s, bitSize) }
func (s *Data) MustUint64(bitSize int) uint64   { return mustUint[uint64](s, bitSize) }
func (s *Data) MustUintptr(bitSize int) uintptr { return mustUint[uintptr](s, bitSize) }
func mustUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr](d *Data, bitSize int) T {
	return assertNextToken[T](d, tokenMustUint[T](bitSize))
}
func tokenMustUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr](bitSize int) func(t Token) T {
	return func(t Token) T { return T(t.MustUint64(bitSize)) }
}
func (s *Data) MustFloat32(bitSize int) float32 { return mustFloat[float32](s, bitSize) }
func (s *Data) MustFloat64(bitSize int) float64 { return mustFloat[float64](s, bitSize) }
func mustFloat[T ~float32 | ~float64](d *Data, bitSize int) T {
	return assertNextToken[T](d, tokenMustFloat[T](bitSize))
}
func tokenMustFloat[T ~float32 | ~float64](bitSize int) func(t Token) T {
	return func(t Token) T { return T(t.MustFloat64(bitSize)) }
}
func (s *Data) MustNumber() string { return assertNextToken[string](s, Token.MustNumber) }
func (s *Data) MustString() string { return assertNextToken[string](s, Token.MustString) }
func (s *Data) MustObject() {
	assertNextToken[struct{}](s, func(t Token) struct{} {
		t.MustObj()
		return struct{}{}
	})
}
func (s *Data) MustArray() {
	assertNextToken[struct{}](s, func(t Token) struct{} {
		t.MustArr()
		return struct{}{}
	})
}
func (s *Data) MustEnd() {
	assertNextToken[struct{}](s, func(t Token) struct{} {
		t.MustEnd()
		return struct{}{}
	})
}
func assertNextToken[T any](d *Data, f func(t Token) T) (v T) {
	if t, ok := d.Peek(); ok {
		v = f(t)
		d.Skip()
	}
	return
}
