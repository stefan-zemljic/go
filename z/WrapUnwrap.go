package z

import (
	"unsafe"
)

func Wrap(bs []byte) string {
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}

func Unwrap(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
