package jso

import (
	"fmt"
)

func (t Token) AddTo(w *Writer) {
	if t == (Token{End}) {
		state := w.State
		switch {
		case state.InObject():
			w.Append('}')
		case state.InArray():
			w.Append(']')
		default:
			panic("cannot end because neither in object nor in array")
		}
		w.PopState()
		return
	}
	switch v := t.data.(type) {
	case nil:
		w.BeforeValue(false)
		w.AppendString("null")
	case bool:
		w.BeforeValue(false)
		if v {
			w.AppendString("true")
		} else {
			w.AppendString("false")
		}
	case number:
		w.BeforeValue(false)
		w.AppendString(string(v))
	case string:
		w.BeforeValue(true)
		w.JsonEscape(v)
	case TokenKind:
		switch v {
		case Array:
			w.BeforeValue(false)
			w.PushState(AtArrayStart)
			w.Append('[')
		case Object:
			w.BeforeValue(false)
			w.PushState(AtObjectStart)
			w.Append('{')
		default:
			panic(fmt.Sprintf("invalid token %v", t))
		}
		return
	}
}
