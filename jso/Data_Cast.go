package jso

func (s *Data) Null() bool {
	_, ok := nextToken[struct{}](s, func(t Token) (struct{}, bool) {
		return struct{}{}, t.Null()
	})
	return ok
}
func (s *Data) Bool() (bool, bool) {
	return nextToken[bool](s, Token.Bool)
}
func (s *Data) Int(bitSize int) (int, bool)     { return tryInt[int](s, bitSize) }
func (s *Data) Int8(bitSize int) (int8, bool)   { return tryInt[int8](s, bitSize) }
func (s *Data) Int16(bitSize int) (int16, bool) { return tryInt[int16](s, bitSize) }
func (s *Data) Int32(bitSize int) (int32, bool) { return tryInt[int32](s, bitSize) }
func (s *Data) Int64(bitSize int) (int64, bool) { return tryInt[int64](s, bitSize) }
func tryInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](d *Data, bitSize int) (T, bool) {
	return nextToken[T](d, tokenTryInt[T](bitSize))
}
func tokenTryInt[T ~int | ~int8 | ~int16 | ~int32 | ~int64](bitSize int) func(t Token) (T, bool) {
	return func(t Token) (T, bool) {
		v, ok := t.Int64(bitSize)
		return T(v), ok
	}
}
func (s *Data) Uint(bitSize int) (uint, bool)       { return tryUint[uint](s, bitSize) }
func (s *Data) Uint8(bitSize int) (uint8, bool)     { return tryUint[uint8](s, bitSize) }
func (s *Data) Uint16(bitSize int) (uint16, bool)   { return tryUint[uint16](s, bitSize) }
func (s *Data) Uint32(bitSize int) (uint32, bool)   { return tryUint[uint32](s, bitSize) }
func (s *Data) Uint64(bitSize int) (uint64, bool)   { return tryUint[uint64](s, bitSize) }
func (s *Data) Uintptr(bitSize int) (uintptr, bool) { return tryUint[uintptr](s, bitSize) }
func tryUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr](d *Data, bitSize int) (T, bool) {
	return nextToken[T](d, tokenTryUint[T](bitSize))
}
func tokenTryUint[T ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr](bitSize int) func(t Token) (T, bool) {
	return func(t Token) (T, bool) {
		v, ok := t.Uint64(bitSize)
		return T(v), ok
	}
}
func (s *Data) Float32(bitSize int) (float32, bool) { return tryFloat[float32](s, bitSize) }
func (s *Data) Float64(bitSize int) (float64, bool) { return tryFloat[float64](s, bitSize) }
func tryFloat[T ~float32 | ~float64](d *Data, bitSize int) (T, bool) {
	return nextToken[T](d, tokenTryFloat[T](bitSize))
}
func tokenTryFloat[T ~float32 | ~float64](bitSize int) func(t Token) (T, bool) {
	return func(t Token) (T, bool) {
		v, ok := t.Float64(bitSize)
		return T(v), ok
	}
}
func (s *Data) Number() (string, bool) { return nextToken[string](s, Token.Number) }
func (s *Data) String() (string, bool) { return nextToken[string](s, Token.String) }
func (s *Data) Obj() bool {
	_, ok := nextToken[struct{}](s, func(t Token) (struct{}, bool) {
		return struct{}{}, t.Obj()
	})
	return ok
}
func (s *Data) Arr() bool {
	_, ok := nextToken[struct{}](s, func(t Token) (struct{}, bool) {
		return struct{}{}, t.Arr()
	})
	return ok
}
func (s *Data) End() bool {
	_, ok := nextToken[struct{}](s, func(t Token) (struct{}, bool) {
		return struct{}{}, t.End()
	})
	return ok
}
func nextToken[T any](d *Data, f func(t Token) (T, bool)) (v T, ok bool) {
	var t Token
	if t, ok = d.Peek(); ok {
		if v, ok = f(t); ok {
			d.Skip()
		}
	}
	return
}
