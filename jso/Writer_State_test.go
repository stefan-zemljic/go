package jso

import (
	"testing"
)

func TestWriter_BeforeValue(t *testing.T) {
	tests := []struct {
		name        string
		initial     WriterState
		valueIsStr  bool
		expectState WriterState
		expectData  string
		expectPanic bool
	}{
		{"AtStart -> AtEnd", AtStart, false, AtEnd, "", false},
		{"AtObjectStart ok", AtObjectStart, true, AfterObjectKey, "", false},
		{"AtObjectStart panic if not string", AtObjectStart, false, 0, "", true},
		{"AfterObjectKey adds colon", AfterObjectKey, true, AfterObjectValue, ":", false},
		{"AfterObjectValue ok (string)", AfterObjectValue, true, AfterObjectKey, ",", false},
		{"AfterObjectValue panic if not string", AfterObjectValue, false, 0, "", true},
		{"AtArrayStart -> AfterArrayValue", AtArrayStart, false, AfterArrayValue, "", false},
		{"AfterArrayValue adds comma", AfterArrayValue, false, AfterArrayValue, ",", false},
		{"AtEnd panic", AtEnd, false, 0, "", true},
		{"Invalid state panic", WriterState(99), false, 0, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &Writer{State: tt.initial}
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("panic = %v, expectPanic %v", r, tt.expectPanic)
				}
			}()
			w.BeforeValue(tt.valueIsStr)
			if !tt.expectPanic {
				if w.State != tt.expectState {
					t.Errorf("expected state %v, got %v", tt.expectState, w.State)
				}
				if string(w.Buffer) != tt.expectData {
					t.Errorf("expected data %q, got %q", tt.expectData, string(w.Buffer))
				}
			}
		})
	}
}
func TestWriter_PushState(t *testing.T) {
	w := &Writer{State: AtStart}
	w.PushState(AtObjectStart)
	if len(w.PrevStates) != 1 || w.PrevStates[0] != AtStart {
		t.Errorf("PrevStates not updated correctly: %+v", w.PrevStates)
	}
	if w.State != AtObjectStart {
		t.Errorf("expected state %v, got %v", AtObjectStart, w.State)
	}
}
func TestWriter_PopState(t *testing.T) {
	t.Run("panic if empty", func(t *testing.T) {
		w := &Writer{}
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("expected panic, got none")
			}
		}()
		w.PopState()
	})
	t.Run("restore previous state", func(t *testing.T) {
		w := &Writer{State: AtEnd, PrevStates: []WriterState{AtStart, AfterObjectKey}}
		w.PopState()
		if w.State != AfterObjectKey {
			t.Errorf("expected state AfterObjectKey, got %v", w.State)
		}
		if len(w.PrevStates) != 1 || w.PrevStates[0] != AtStart {
			t.Errorf("PrevStates not shortened correctly: %+v", w.PrevStates)
		}
	})
}
