package jso

import (
	"reflect"
	"testing"
)

type mockAdapter struct {
	writeReturn bool
	readValue   any
	readOk      bool
}

func (m *mockAdapter) Write(any, Stream) bool {
	return m.writeReturn
}
func (m *mockAdapter) Read(reflect.Type, Stream) (any, bool) {
	return m.readValue, m.readOk
}
func TestRegistry_Write(t *testing.T) {
	type sample struct{}
	tests := []struct {
		name        string
		setup       func() *Registry
		expectPanic bool
	}{
		{
			name: "write succeeds via type adapter",
			setup: func() *Registry {
				var r Registry
				r.Adapter(reflect.TypeOf(sample{}), &mockAdapter{writeReturn: true})
				return &r
			},
		},
		{
			name: "write succeeds via fallback adapter",
			setup: func() *Registry {
				var r Registry
				r.Adapter(reflect.TypeOf(sample{}), &mockAdapter{writeReturn: false})
				r.Fallback(&mockAdapter{writeReturn: true})
				return &r
			},
		},
		{
			name: "write panics when no adapter works",
			setup: func() *Registry {
				var r Registry
				r.Adapter(reflect.TypeOf(sample{}), &mockAdapter{writeReturn: false})
				r.Fallback(&mockAdapter{writeReturn: false})
				return &r
			},
			expectPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := tt.setup()
			stream := Stream{}
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("panic = %v, expectPanic %v", r, tt.expectPanic)
				}
			}()
			reg.Write(sample{}, stream)
		})
	}
}
func TestRegistry_Read(t *testing.T) {
	type sample struct{}
	tests := []struct {
		name        string
		setup       func() *Registry
		expectValue any
		expectPanic bool
	}{
		{
			name: "read succeeds via type adapter",
			setup: func() *Registry {
				var r Registry
				r.Adapter(reflect.TypeOf(sample{}), &mockAdapter{readValue: 123, readOk: true})
				return &r
			},
			expectValue: 123,
		},
		{
			name: "read succeeds via fallback adapter",
			setup: func() *Registry {
				var r Registry
				r.Fallback(&mockAdapter{readValue: "ok", readOk: true})
				return &r
			},
			expectValue: "ok",
		},
		{
			name: "read panics when no adapter works",
			setup: func() *Registry {
				var r Registry
				r.Adapter(reflect.TypeOf(sample{}), &mockAdapter{readOk: false})
				r.Fallback(&mockAdapter{readOk: false})
				return &r
			},
			expectPanic: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reg := tt.setup()
			stream := Stream{}
			defer func() {
				if r := recover(); (r != nil) != tt.expectPanic {
					t.Errorf("panic = %v, expectPanic %v", r, tt.expectPanic)
				}
			}()
			if !tt.expectPanic {
				val := reg.Read(reflect.TypeOf(sample{}), stream)
				if val != tt.expectValue {
					t.Errorf("expected %v, got %v", tt.expectValue, val)
				}
			} else {
				_ = reg.Read(reflect.TypeOf(sample{}), stream)
			}
		})
	}
}
