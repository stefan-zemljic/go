package jso

import (
	"fmt"
	"slices"
	"unicode/utf8"
)

func Escape(b []byte, str string) []byte {
	if len(str) == 0 {
		return append(b, '"', '"')
	}
	b = append(b, '"')
	for _, r := range str {
		switch r {
		case '\\':
			b = append(b, '\\', '\\')
		case '"':
			b = append(b, '\\', '"')
		case '\b':
			b = append(b, '\\', 'b')
		case '\t':
			b = append(b, '\\', 't')
		case '\n':
			b = append(b, '\\', 'n')
		case '\f':
			b = append(b, '\\', 'f')
		case '\r':
			b = append(b, '\\', 'r')
		default:
			if r < 0x20 {
				b = fmt.Appendf(b, "\\u%04x", r)
			} else if r > 0xFFFF {
				r1 := 0xD800 + ((r - 0x10000) >> 10)
				r2 := 0xDC00 + (r & 0x3FF)
				b = fmt.Appendf(b, "\\u%04x\\u%04x", r1, r2)
			} else {
				b = slices.Grow(b, utf8.UTFMax)
				n := utf8.EncodeRune(b[len(b):len(b)+utf8.UTFMax], r)
				b = b[:len(b)+n]
			}
		}
	}
	return append(b, '"')
}
