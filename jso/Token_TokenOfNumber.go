package jso

import (
	"fmt"
	"regexp"
)

var jsonNumberRegexp = regexp.MustCompile(`^-?(0|[1-9]\d*)(\.\d+)?([eE][+-]?\d+)?$`)

func TokenOfNumber(s string) Token {
	if !jsonNumberRegexp.MatchString(s) {
		panic(fmt.Sprintf("invalid JSON number: %q", s))
	}
	return Token{number(s)}
}
