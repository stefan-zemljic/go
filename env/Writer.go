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

type Writer struct {
	buf []byte
}

func (s *Writer) Set(name, value string) {
	if s.buf == nil {
		s.buf = append(s.buf, hex2str...)
	}
	_, _ = fmt.Fprintf(s, "export %s=$(hex2str %x)\n", name, []byte(value))
}

func (s *Writer) Unset(name string) {
	if s.buf == nil {
		s.buf = append(s.buf, hex2str...)
	}
	_, _ = fmt.Fprintf(s, "unset %s\n", name)
}

func (s *Writer) Bytes() []byte { return s.buf }

func (s *Writer) String() string { return string(s.buf) }

func (s *Writer) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(s.buf)
	return int64(n), err
}

func (s *Writer) Reset() {
	s.buf = nil
}

func (s *Writer) Write(bs []byte) (int, error) {
	s.buf = append(s.buf, bs...)
	return len(bs), nil
}
