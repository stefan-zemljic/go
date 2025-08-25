package jso

import (
	"fmt"
	"strconv"
)

func TokenOf(v any) Token {
	switch v := v.(type) {
	case nil:
		return Token{nil}
	case bool:
		return Token{v}
	case int:
		return Token{number(strconv.FormatInt(int64(v), 10))}
	case int8:
		return Token{number(strconv.FormatInt(int64(v), 10))}
	case int16:
		return Token{number(strconv.FormatInt(int64(v), 10))}
	case int32:
		return Token{number(strconv.FormatInt(int64(v), 10))}
	case int64:
		return Token{number(strconv.FormatInt(v, 10))}
	case uint:
		return Token{number(strconv.FormatUint(uint64(v), 10))}
	case uint8:
		return Token{number(strconv.FormatUint(uint64(v), 10))}
	case uint16:
		return Token{number(strconv.FormatUint(uint64(v), 10))}
	case uint32:
		return Token{number(strconv.FormatUint(uint64(v), 10))}
	case uint64:
		return Token{number(strconv.FormatUint(v, 10))}
	case float32:
		return Token{number(strconv.FormatFloat(float64(v), 'g', -1, 32))}
	case float64:
		return Token{number(strconv.FormatFloat(v, 'g', -1, 64))}
	case string:
		return Token{v}
	case TokenKind:
		switch v {
		case Array, Object, End:
			return Token{data: v}
		default:
			panic(fmt.Sprintf("token kind %d cannot be used as value", v))
		}
	default:
		panic(fmt.Sprintf("unsupported token type %T", v))
	}
}
