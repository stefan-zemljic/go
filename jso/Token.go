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
	Arr
	Obj
	End
)

func (t TokenKind) String() string {
	switch t {
	case Null:
		return "Null"
	case Bool:
		return "Bool"
	case Number:
		return "Number"
	case String:
		return "String"
	case Arr:
		return "Arr"
	case Obj:
		return "Obj"
	case End:
		return "End"
	default:
		panic(fmt.Sprintf("invalid TokenKind %d", t))
	}
}

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
		case Arr:
			return Arr
		case Obj:
			return Obj
		case End:
			return End
		default:
		}
	}
	panic(fmt.Sprintf("invalid token %v", t.data))
}
