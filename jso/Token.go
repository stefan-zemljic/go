package jso

import (
	"fmt"
)

type Token struct {
	data any
}
type TokenKind int

const (
	Null TokenKind = iota
	Bool
	Number
	String
	Array
	Object
	End
)

type number string

func (t Token) Kind() TokenKind {
	switch c := t.data.(type) {
	case nil:
		return Null
	case bool:
		return Bool
	case number:
		return Number
	case string:
		return String
	case TokenKind:
		switch c {
		case Array:
			return Array
		case Object:
			return Object
		case End:
			return End
		default:
		}
	}
	panic(fmt.Sprintf("invalid token %v", t.data))
}
