package env

import (
	_ "embed"
	"fmt"
	"io"
)

const hex2str = `hex2str() {
  local hex="$1"
  local esc=""
  local i
  for ((i=0; i<${#hex}; i+=2)); do
    esc+="\\x${hex:$i:2}"
  done
  eval "echo \$'$esc'"
}
`

type Emitter struct {
	buf []byte
}

func (e *Emitter) Set(name, value string) {
	if e.buf == nil {
		e.buf = append(e.buf, hex2str...)
	}
	_, _ = fmt.Fprintf(e, "export %s=$(hex2str %x)\n", name, []byte(value))
}
func (e *Emitter) Unset(name string) {
	if e.buf == nil {
		e.buf = append(e.buf, hex2str...)
	}
	_, _ = fmt.Fprintf(e, "unset %s\n", name)
}
func (e *Emitter) Bytes() []byte  { return e.buf }
func (e *Emitter) String() string { return string(e.buf) }
func (e *Emitter) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(e.buf)
	return int64(n), err
}
func (e *Emitter) Reset() {
	e.buf = nil
}
func (e *Emitter) Write(bs []byte) (int, error) {
	e.buf = append(e.buf, bs...)
	return len(bs), nil
}
