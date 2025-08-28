package jso

import (
	"bytes"
	"encoding/json"
	"io"
	"strconv"
	"unsafe"
)

type Writer struct {
	buf    []byte
	state  state
	states []state
}
type state int

const (
	atStart state = iota
	atObjectStart
	afterObjectKey
	afterObjectValue
	atArrayStart
	afterArrayValue
	atEnd
)

func (w *Writer) Null()             { w.beforeValue(); w.buf = append(w.buf, "null"...) }
func (w *Writer) Bool(v bool)       { w.beforeValue(); w.buf = strconv.AppendBool(w.buf, v) }
func (w *Writer) Int(v int)         { w.beforeValue(); w.buf = strconv.AppendInt(w.buf, int64(v), 10) }
func (w *Writer) Int8(v int8)       { w.beforeValue(); w.buf = strconv.AppendInt(w.buf, int64(v), 10) }
func (w *Writer) Int16(v int16)     { w.beforeValue(); w.buf = strconv.AppendInt(w.buf, int64(v), 10) }
func (w *Writer) Int32(v int32)     { w.beforeValue(); w.buf = strconv.AppendInt(w.buf, int64(v), 10) }
func (w *Writer) Int64(v int64)     { w.beforeValue(); w.buf = strconv.AppendInt(w.buf, v, 10) }
func (w *Writer) Uint(v uint)       { w.beforeValue(); w.buf = strconv.AppendUint(w.buf, uint64(v), 10) }
func (w *Writer) Uint8(v uint8)     { w.beforeValue(); w.buf = strconv.AppendUint(w.buf, uint64(v), 10) }
func (w *Writer) Uint16(v uint16)   { w.beforeValue(); w.buf = strconv.AppendUint(w.buf, uint64(v), 10) }
func (w *Writer) Uint32(v uint32)   { w.beforeValue(); w.buf = strconv.AppendUint(w.buf, uint64(v), 10) }
func (w *Writer) Uint64(v uint64)   { w.beforeValue(); w.buf = strconv.AppendUint(w.buf, v, 10) }
func (w *Writer) Uintptr(v uintptr) { w.beforeValue(); w.buf = strconv.AppendUint(w.buf, uint64(v), 10) }
func (w *Writer) Float32(v float32) {
	w.beforeValue()
	w.buf = strconv.AppendFloat(w.buf, float64(v), 'g', -1, 32)
}
func (w *Writer) Float64(v float64) {
	w.beforeValue()
	w.buf = strconv.AppendFloat(w.buf, v, 'g', -1, 64)
}
func (w *Writer) Number(v string) { w.beforeValue(); w.buf = append(w.buf, v...) }
func (w *Writer) String(v string) { w.beforeString(); w.buf = Escape(w.buf, v) }
func (w *Writer) Array()          { w.beforeValue(); w.buf = append(w.buf, '['); w.pushState(atArrayStart) }
func (w *Writer) Object()         { w.beforeValue(); w.buf = append(w.buf, '{'); w.pushState(atObjectStart) }
func (w *Writer) End() {
	if len(w.states) == 0 {
		panic("jso: cannot end, no open array or object")
	}
	var end byte
	switch w.state {
	case atStart, atEnd:
		panic("jso: cannot end, not in array or object")
	case atObjectStart, afterObjectValue:
		end = '}'
	case atArrayStart, afterArrayValue:
		end = ']'
	case afterObjectKey:
		panic("jso: cannot end after object key")
	}
	lastIndex := len(w.states) - 1
	w.state, w.states = w.states[lastIndex], w.states[:lastIndex]
	w.buf = append(w.buf, end)
}
func (w *Writer) Compact() string {
	return unsafe.String(unsafe.SliceData(w.buf), len(w.buf))
}
func (w *Writer) CompactBytes() []byte {
	return w.buf
}
func (w *Writer) CompactTo(out io.Writer) {
	_, err := out.Write(w.buf)
	if err != nil {
		panic(err)
	}
}
func (w *Writer) Pretty() string {
	return w.Indent("", "  ")
}
func (w *Writer) PrettyBytes() []byte {
	var buf bytes.Buffer
	err := json.Indent(&buf, w.buf, "", "  ")
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
func (w *Writer) PrettyTo(out io.Writer) {
	_, err := out.Write(w.PrettyBytes())
	if err != nil {
		panic(err)
	}
}
func (w *Writer) Indent(prefix, indent string) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, w.buf, prefix, indent)
	if err != nil {
		panic(err)
	}
	return buf.String()
}
func (w *Writer) IndentBytes(prefix, indent string) []byte {
	var buf bytes.Buffer
	err := json.Indent(&buf, w.buf, prefix, indent)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}
func (w *Writer) IndentTo(out io.Writer, prefix, indent string) {
	_, err := out.Write(w.IndentBytes(prefix, indent))
	if err != nil {
		panic(err)
	}
}
func (w *Writer) beforeValue() {
	switch w.state {
	case atStart:
		w.state = atEnd
	case atObjectStart:
		panic("jso: object keys must be strings")
	case afterObjectKey:
		w.buf = append(w.buf, ':')
		w.state = afterObjectValue
	case afterObjectValue:
		panic("jso: object keys must be strings")
	case atArrayStart:
		w.state = afterArrayValue
	case afterArrayValue:
		w.buf = append(w.buf, ',')
	case atEnd:
		panic("jso: cannot write value, document already complete")
	}
}
func (w *Writer) beforeString() {
	switch w.state {
	case atObjectStart:
		w.state = afterObjectKey
	case afterObjectValue:
		w.buf = append(w.buf, ',')
		w.state = afterObjectKey
	default:
		w.beforeValue()
	}
}
func (w *Writer) pushState(s state) {
	w.states = append(w.states, w.state)
	w.state = s
}
