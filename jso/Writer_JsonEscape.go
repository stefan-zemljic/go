package jso

import (
	"fmt"
	"unicode/utf8"
)

func (s *Writer) JsonEscape(str string) {
	if len(str) == 0 {
		s.AppendString(`""`)
		return
	}
	s.Append('"')
	for _, r := range str {
		switch r {
		case '\\':
			s.AppendString(`\\`)
		case '"':
			s.AppendString(`\"`)
		case '\b':
			s.AppendString(`\b`)
		case '\t':
			s.AppendString(`\t`)
		case '\n':
			s.AppendString(`\n`)
		case '\f':
			s.AppendString(`\f`)
		case '\r':
			s.AppendString(`\r`)
		default:
			if r < 0x20 {
				_, _ = fmt.Fprintf(s, "\\u%04x", r)
			} else if r > 0xFFFF {
				r1 := 0xD800 + ((r - 0x10000) >> 10)
				r2 := 0xDC00 + (r & 0x3FF)
				_, _ = fmt.Fprintf(s, "\\u%04x\\u%04x", r1, r2)
			} else {
				var w [utf8.UTFMax]byte
				n := utf8.EncodeRune(w[:], r)
				_, _ = s.Write(w[:n])
			}
		}
	}
	s.Append('"')
}
