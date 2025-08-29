package jso

import (
	"regexp"
	"strconv"
)

type Reader struct {
	buf  []byte
	has  bool
	next any
}

type num []byte

func (s *Reader) Save() []byte        { return s.buf }
func (s *Reader) Restore(save []byte) { s.buf = save }

var jsonNumberRegexp = regexp.MustCompile(`^-?(0|[1-9][0-9]*)(\.[0-9]+)?([eE][+-]?[0-9]+)?`)
var jsonStringRegexp = regexp.MustCompile(`"(?:\\["\\/bfnrt]|\\u[0-9a-fA-F]{4}|[^"\\\x00-\x1F])*"`)

func (s *Reader) Kind() Kind {
	if !s.has {
		s.prepareNext()
		if !s.has {
			return EOF
		}
	}
	switch v := s.next.(type) {
	case nil:
		return Null
	case bool:
		return Bool
	case num:
		return Number
	case []byte:
		return String
	case byte:
		switch v {
		case '[':
			return Array
		case '{':
			return Object
		case ']', '}':
			return End
		}
	}
	panic("jso: invalid state")
}

func (s *Reader) Null() bool {
	if !s.has {
		s.prepareNext()
		if !s.has {
			return false
		}
	}
	if s.next == nil {
		s.has = false
		s.next = nil
		return true
	}
	return false
}

func (s *Reader) MustNull() {
	if !s.has {
		s.prepareNext()
		if !s.has {
			panic("jso: unexpected end of input")
		}
	}
	if s.next != nil {
		panic("jso: expected null")
	}
	s.has = false
	s.next = nil
}

func (s *Reader) Bool() (bool, bool) { return read[bool](s) }
func (s *Reader) MustBool() bool     { return mustRead[bool](s) }

func (s *Reader) Int() (int, bool)     { return readInt[int](s, 0) }
func (s *Reader) MustInt() int         { return mustReadInt[int](s, 0) }
func (s *Reader) Int8() (int8, bool)   { return readInt[int8](s, 8) }
func (s *Reader) MustInt8() int8       { return mustReadInt[int8](s, 8) }
func (s *Reader) Int16() (int16, bool) { return readInt[int16](s, 16) }
func (s *Reader) MustInt16() int16     { return mustReadInt[int16](s, 16) }
func (s *Reader) Int32() (int32, bool) { return readInt[int32](s, 32) }
func (s *Reader) MustInt32() int32     { return mustReadInt[int32](s, 32) }
func (s *Reader) Int64() (int64, bool) { return readInt[int64](s, 64) }
func (s *Reader) MustInt64() int64     { return mustReadInt[int64](s, 64) }

func (s *Reader) Uint() (uint, bool)       { return readUint[uint](s, 0) }
func (s *Reader) MustUint() uint           { return mustReadUint[uint](s, 0) }
func (s *Reader) Uint8() (uint8, bool)     { return readUint[uint8](s, 8) }
func (s *Reader) MustUint8() uint8         { return mustReadUint[uint8](s, 8) }
func (s *Reader) Uint16() (uint16, bool)   { return readUint[uint16](s, 16) }
func (s *Reader) MustUint16() uint16       { return mustReadUint[uint16](s, 16) }
func (s *Reader) Uint32() (uint32, bool)   { return readUint[uint32](s, 32) }
func (s *Reader) MustUint32() uint32       { return mustReadUint[uint32](s, 32) }
func (s *Reader) Uint64() (uint64, bool)   { return readUint[uint64](s, 64) }
func (s *Reader) MustUint64() uint64       { return mustReadUint[uint64](s, 64) }
func (s *Reader) Uintptr() (uintptr, bool) { return readUint[uintptr](s, strconv.IntSize) }
func (s *Reader) MustUintptr() uintptr     { return mustReadUint[uintptr](s, strconv.IntSize) }

func (s *Reader) Float32() (float32, bool) {
	if v, ok := s.Float64(); ok {
		return float32(v), true
	}
	return 0, false
}
func (s *Reader) MustFloat32() float32 {
	return float32(s.MustFloat64())
}
func (s *Reader) Float64() (float64, bool) {
	v, ok := read[num](s)
	if !ok {
		return 0, false
	}
	n, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		return 0, false
	}
	return n, true
}
func (s *Reader) MustFloat64() float64 {
	v := mustRead[num](s)
	n, err := strconv.ParseFloat(string(v), 64)
	if err != nil {
		panic("jso: invalid number")
	}
	return n
}

func (s *Reader) Number() (string, bool) {
	if v, ok := read[num](s); ok {
		return string(v), true
	}
	return "", false
}
func (s *Reader) MustNumber() string { return string(mustRead[num](s)) }

func (s *Reader) String() (string, bool) {
	if v, ok := read[[]byte](s); ok {
		return string(v), true
	}
	return "", false
}
func (s *Reader) MustString() string { return string(mustRead[[]byte](s)) }

func (s *Reader) Array() bool  { return s.just('[') }
func (s *Reader) MustArray()   { s.must('[') }
func (s *Reader) Object() bool { return s.just('{') }
func (s *Reader) MustObject()  { s.must('{') }
func (s *Reader) End() bool    { return s.just(']') || s.just('}') }
func (s *Reader) MustEnd() {
	if !s.End() {
		panic("jso: expected end of array or object")
	}
}

func (s *Reader) must(b byte) {
	if !s.just(b) {
		panic("jso: expected different token")
	}
}

func (s *Reader) just(b byte) bool {
	if !s.has {
		s.prepareNext()
		if !s.has {
			return false
		}
	}
	v, ok := s.next.(byte)
	if !ok || v != b {
		return false
	}
	s.has = false
	s.next = nil
	return true
}

func readInt[T int | int8 | int16 | int32 | int64](s *Reader, bitSize int) (T, bool) {
	v, ok := read[num](s)
	if !ok {
		var zero T
		return zero, false
	}
	n, err := strconv.ParseInt(string(v), 10, bitSize)
	if err != nil {
		var zero T
		return zero, false
	}
	return T(n), true
}

func mustReadInt[T int | int8 | int16 | int32 | int64](s *Reader, bitSize int) T {
	v := mustRead[num](s)
	n, err := strconv.ParseInt(string(v), 10, bitSize)
	if err != nil {
		panic("jso: invalid number")
	}
	return T(n)
}

func readUint[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](s *Reader, bitSize int) (T, bool) {
	v, ok := read[num](s)
	if !ok {
		var zero T
		return zero, false
	}
	n, err := strconv.ParseUint(string(v), 10, bitSize)
	if err != nil {
		var zero T
		return zero, false
	}
	return T(n), true
}

func mustReadUint[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](s *Reader, bitSize int) T {
	v := mustRead[num](s)
	n, err := strconv.ParseUint(string(v), 10, bitSize)
	if err != nil {
		panic("jso: invalid number")
	}
	return T(n)
}

func read[T any](s *Reader) (T, bool) {
	if !s.has {
		s.prepareNext()
		if !s.has {
			var zero T
			return zero, false
		}
	}
	v, ok := s.next.(T)
	if !ok {
		var zero T
		return zero, false
	}
	s.has = false
	s.next = nil
	return v, true
}

func mustRead[T any](s *Reader) T {
	if !s.has {
		s.prepareNext()
		if !s.has {
			panic("jso: unexpected end of input")
		}
	}
	v, ok := s.next.(T)
	if !ok {
		panic("jso: expected different token")
	}
	s.has = false
	s.next = nil
	return v
}

func (s *Reader) prepareNext() {
	for {
		if len(s.buf) == 0 {
			s.has = false
			s.next = nil
			return
		}
		switch s.buf[0] {
		case ' ', '\t', '\n', '\r', ',', ':':
			s.buf = s.buf[1:]
			continue
		}
		break
	}
	b := s.buf[0]
	switch b {
	case 'n':
		if len(s.buf) < 4 {
			panic("jso: invalid JSON input")
		}
		s.buf = s.buf[4:]
	case 't':
		if len(s.buf) < 4 {
			panic("jso: invalid JSON input")
		}
		s.buf = s.buf[4:]
		s.next = true
	case 'f':
		if len(s.buf) < 5 {
			panic("jso: invalid JSON input")
		}
		s.buf = s.buf[5:]
		s.next = false
	case '[', '{', ']', '}':
		s.next = b
		s.buf = s.buf[1:]
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		m := jsonNumberRegexp.FindIndex(s.buf)
		if m == nil || m[0] != 0 {
			panic("jso: invalid JSON input")
		}
		end := m[1]
		s.next = num(s.buf[m[0]:end])
		s.buf = s.buf[end:]
	case '"':
		m := jsonStringRegexp.FindIndex(s.buf)
		if m == nil || m[0] != 0 {
			panic("jso: invalid JSON input")
		}
		end := m[1]
		s.next = s.buf[m[0]+1 : end-1]
		s.buf = s.buf[end:]
	default:
		panic("jso: invalid JSON input")
	}
	s.has = true
}
