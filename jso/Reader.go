package jso

type Reader interface {
	Null() bool
	Bool() (bool, bool)
	Int() (int, bool)
	Int8() (int8, bool)
	Int16() (int16, bool)
	Int32() (int32, bool)
	Int64() (int64, bool)
	Uint() (uint, bool)
	Uint8() (uint8, bool)
	Uint16() (uint16, bool)
	Uint32() (uint32, bool)
	Uint64() (uint64, bool)
	Uintptr() (uintptr, bool)
	Float32() (float32, bool)
	Float64() (float64, bool)
	Number() (string, bool)
	String() (string, bool)
	Array() bool
	Object() bool
	End() bool
	MustNull()
	MustBool() bool
	MustInt() int
	MustInt8() int8
	MustInt16() int16
	MustInt32() int32
	MustInt64() int64
	MustUint() uint
	MustUint8() uint8
	MustUint16() uint16
	MustUint32() uint32
	MustUint64() uint64
	MustUintptr() uintptr
	MustFloat32() float32
	MustFloat64() float64
	MustNumber() string
	MustString() string
	MustArray()
	MustObject()
	MustEnd()
}
