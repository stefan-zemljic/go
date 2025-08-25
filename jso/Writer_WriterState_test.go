package jso

import (
	"testing"
)

func TestWriterState_InObject(t *testing.T) {
	tests := []struct {
		state    WriterState
		expected bool
	}{
		{AtStart, false},
		{AtObjectStart, true},
		{AfterObjectKey, true},
		{AfterObjectValue, true},
		{AtArrayStart, false},
		{AfterArrayValue, false},
		{AtEnd, false},
	}
	for _, tt := range tests {
		if got := tt.state.InObject(); got != tt.expected {
			t.Errorf("InObject(%v) = %v, expected %v", tt.state, got, tt.expected)
		}
	}
}
func TestWriterState_InArray(t *testing.T) {
	tests := []struct {
		state    WriterState
		expected bool
	}{
		{AtStart, false},
		{AtObjectStart, false},
		{AfterObjectKey, false},
		{AfterObjectValue, false},
		{AtArrayStart, true},
		{AfterArrayValue, true},
		{AtEnd, false},
	}
	for _, tt := range tests {
		if got := tt.state.InArray(); got != tt.expected {
			t.Errorf("InArray(%v) = %v, expected %v", tt.state, got, tt.expected)
		}
	}
}
