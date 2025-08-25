package jso

import (
	"bytes"
	"encoding/json"
	"io"
	"unsafe"
)

type Writer struct {
	Buffer     []byte
	State      WriterState
	PrevStates []WriterState
}

func (s *Writer) Add(value any) {
	TokenOf(value).AddTo(s)
}
func (s *Writer) AddAll(values ...any) {
	for _, v := range values {
		TokenOf(v).AddTo(s)
	}
}
func (s *Writer) Append(bs ...byte) {
	s.Buffer = append(s.Buffer, bs...)
}
func (s *Writer) AppendString(str string) {
	s.Buffer = append(s.Buffer, str...)
}

var _ io.Writer = (*Writer)(nil)

func (s *Writer) Write(data []byte) (int, error) {
	s.Buffer = append(s.Buffer, data...)
	return len(data), nil
}
func (s *Writer) String() string {
	return unsafe.String(unsafe.SliceData(s.Buffer), len(s.Buffer))
}
func (s *Writer) Pretty() string {
	return s.Indent("", "  ")
}
func (s *Writer) Indent(prefix, indent string) string {
	var outBuf bytes.Buffer
	err := json.Indent(&outBuf, s.Buffer, prefix, indent)
	if err != nil {
		panic(err)
	}
	return outBuf.String()
}
func (s *Writer) Reset() {
	s.Buffer = s.Buffer[:0]
	s.State = AtStart
	s.PrevStates = s.PrevStates[:0]
}
func StringOf(values ...any) string {
	var w Writer
	w.AddAll(values...)
	return w.String()
}
func Pretty(values ...any) string {
	var w Writer
	w.AddAll(values...)
	return w.Pretty()
}
func Indent(prefix, indent string, values ...any) string {
	var w Writer
	w.AddAll(values...)
	return w.Indent(prefix, indent)
}
