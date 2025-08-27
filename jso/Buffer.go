package jso

import (
	"bytes"
	"encoding/json"
	"io"
	"unsafe"
)

type Buffer struct {
	Buffer     []byte
	State      BufferState
	PrevStates []BufferState
}

func (s *Buffer) Add(value any) {
	TokenOf(value).AddTo(s)
}
func (s *Buffer) AddAll(values ...any) {
	for _, v := range values {
		TokenOf(v).AddTo(s)
	}
}
func (s *Buffer) Append(bs ...byte) {
	s.Buffer = append(s.Buffer, bs...)
}
func (s *Buffer) AppendString(str string) {
	s.Buffer = append(s.Buffer, str...)
}

var _ io.Writer = (*Buffer)(nil)

func (s *Buffer) Write(data []byte) (int, error) {
	s.Buffer = append(s.Buffer, data...)
	return len(data), nil
}
func (s *Buffer) Json() string {
	return unsafe.String(unsafe.SliceData(s.Buffer), len(s.Buffer))
}
func (s *Buffer) Pretty() string {
	return s.Indent("", "  ")
}
func (s *Buffer) Indent(prefix, indent string) string {
	var outBuf bytes.Buffer
	err := json.Indent(&outBuf, s.Buffer, prefix, indent)
	if err != nil {
		panic(err)
	}
	return outBuf.String()
}
func (s *Buffer) Reset() {
	s.Buffer = s.Buffer[:0]
	s.State = AtStart
	s.PrevStates = s.PrevStates[:0]
}
func (s *Buffer) Obj() { s.Add(Obj) }
func (s *Buffer) Arr() { s.Add(Arr) }
func (s *Buffer) End() { s.Add(End) }
func (s *Buffer) Nil() { s.Add(nil) }
func Json(values ...any) string {
	var b Buffer
	b.AddAll(values...)
	return b.Json()
}
func Pretty(values ...any) string {
	var b Buffer
	b.AddAll(values...)
	return b.Pretty()
}
func Indent(prefix, indent string, values ...any) string {
	var b Buffer
	b.AddAll(values...)
	return b.Indent(prefix, indent)
}
