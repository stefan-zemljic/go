package jso

import (
	"strconv"
	"strings"
)

func Unescape(bs []byte) string {
	bs = bs[1 : len(bs)-1]
	var sb strings.Builder
	sb.Grow(len(bs))
	i := 0
	for i < len(bs) {
		b := bs[i]
		switch b {
		case '\\':
			if i+1 >= len(bs) {
				panic("jso: invalid escape sequence")
			}
			switch bs[i+1] {
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			case '/':
				sb.WriteByte('/')
			case 'b':
				sb.WriteByte('\b')
			case 'f':
				sb.WriteByte('\f')
			case 'n':
				sb.WriteByte('\n')
			case 'r':
				sb.WriteByte('\r')
			case 't':
				sb.WriteByte('\t')
			case 'u':
				if i+6 > len(bs) {
					panic("jso: invalid unicode escape sequence")
				}
				hex := string(bs[i+2 : i+6])
				v, err := strconv.ParseUint(hex, 16, 16)
				if err != nil {
					panic("jso: invalid unicode escape sequence")
				}
				sb.WriteRune(rune(v))
				i += 4
			default:
				panic("jso: invalid escape sequence")
			}
			i += 2
		default:
			sb.WriteByte(b)
			i++
		}
	}
	return sb.String()
}
