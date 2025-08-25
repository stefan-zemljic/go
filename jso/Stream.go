package jso

type Stream struct {
	Tokens []Token
}
type Num string

func (s *Stream) Add(value any) {
	s.Tokens = append(s.Tokens, TokenOf(value))
}
func (s *Stream) AddAll(values ...any) {
	for _, value := range values {
		s.Tokens = append(s.Tokens, TokenOf(value))
	}
}
func (s *Stream) String() string {
	var w Writer
	for _, t := range s.Tokens {
		t.AddTo(&w)
	}
	return w.String()
}
func (s *Stream) Pretty() string {
	return s.Indent("", "  ")
}
func (s *Stream) Indent(prefix, indent string) string {
	var w Writer
	for _, t := range s.Tokens {
		t.AddTo(&w)
	}
	return w.Indent(prefix, indent)
}
func (s *Stream) Peek() (Token, bool) {
	if len(s.Tokens) == 0 {
		return Token{}, false
	}
	return s.Tokens[0], true
}
func (s *Stream) Read() Token {
	if len(s.Tokens) == 0 {
		panic("no more tokens to read")
	}
	t := s.Tokens[0]
	s.Tokens = s.Tokens[1:]
	return t
}
func (s *Stream) IsEmpty() bool {
	return len(s.Tokens) == 0
}
func (s *Stream) More() bool {
	return len(s.Tokens) > 0
}
