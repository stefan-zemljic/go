package jso

import (
	"fmt"
)

func (t Token) AddTo(b *Buffer) {
	if t == (Token{End}) {
		state := b.State
		switch {
		case state.InObject():
			b.Append('}')
		case state.InArray():
			b.Append(']')
		default:
			panic("cannot end because neither in object nor in array")
		}
		b.PopState()
		return
	}
	switch v := t.data.(type) {
	case nil:
		b.BeforeValue(false)
		b.AppendString("null")
	case bool:
		b.BeforeValue(false)
		if v {
			b.AppendString("true")
		} else {
			b.AppendString("false")
		}
	case number:
		b.BeforeValue(false)
		b.AppendString(string(v))
	case string:
		b.BeforeValue(true)
		b.JsonEscape(v)
	case TokenKind:
		switch v {
		case Arr:
			b.BeforeValue(false)
			b.PushState(AtArrayStart)
			b.Append('[')
		case Obj:
			b.BeforeValue(false)
			b.PushState(AtObjectStart)
			b.Append('{')
		default:
			panic(fmt.Sprintf("invalid token %v", t))
		}
		return
	}
}
