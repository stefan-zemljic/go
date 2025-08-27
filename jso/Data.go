package jso

type Data struct {
	Tokens []Token
}

func (s *Data) Add(value any) {
	s.Tokens = append(s.Tokens, TokenOf(value))
}
func (s *Data) AddAll(values ...any) {
	for _, value := range values {
		s.Tokens = append(s.Tokens, TokenOf(value))
	}
}
func (s *Data) Json() string {
	var b Buffer
	for _, t := range s.Tokens {
		t.AddTo(&b)
	}
	return b.Json()
}
func (s *Data) Pretty() string {
	return s.Indent("", "  ")
}
func (s *Data) Indent(prefix, indent string) string {
	var b Buffer
	for _, t := range s.Tokens {
		t.AddTo(&b)
	}
	return b.Indent(prefix, indent)
}
func (s *Data) Peek() (Token, bool) {
	if len(s.Tokens) == 0 {
		return Token{}, false
	}
	return s.Tokens[0], true
}
func (s *Data) Skip() {
	if len(s.Tokens) == 0 {
		panic("no token to skip")
	}
	s.Tokens = s.Tokens[1:]
}
func (s *Data) Read() Token {
	if len(s.Tokens) == 0 {
		panic("no token to read")
	}
	t := s.Tokens[0]
	s.Tokens = s.Tokens[1:]
	return t
}
func (s *Data) IsEmpty() bool {
	return len(s.Tokens) == 0
}
func (s *Data) More() bool {
	return len(s.Tokens) > 0
}
