package jso

import (
	"fmt"
)

func (s *Buffer) BeforeValue(valueIsString bool) {
	switch s.State {
	case AtStart:
		s.State = AtEnd
	case AtObjectStart:
		if !valueIsString {
			panic("object key must be string")
		}
		s.State = AfterObjectKey
	case AfterObjectKey:
		s.Append(':')
		s.State = AfterObjectValue
	case AfterObjectValue:
		if !valueIsString {
			panic("object key must be string")
		}
		s.Append(',')
		s.State = AfterObjectKey
	case AtArrayStart:
		s.State = AfterArrayValue
	case AfterArrayValue:
		s.Append(',')
		s.State = AfterArrayValue
	case AtEnd:
		panic("cannot add value after end")
	default:
		panic(fmt.Sprintf("invalid state %d", s.State))
	}
}
func (s *Buffer) PushState(state BufferState) {
	s.PrevStates = append(s.PrevStates, s.State)
	s.State = state
}
func (s *Buffer) PopState() {
	if len(s.PrevStates) == 0 {
		panic("cannot pop state because no previous state")
	}
	s.State = s.PrevStates[len(s.PrevStates)-1]
	s.PrevStates = s.PrevStates[:len(s.PrevStates)-1]
}
