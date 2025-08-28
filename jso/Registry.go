package jso

import (
	"reflect"
)

type Registry struct {
	adapters map[reflect.Type]Adapter
}

func (s *Registry) Register(t reflect.Type, adapter Adapter) {
	if s.adapters == nil {
		s.adapters = make(map[reflect.Type]Adapter)
	}
	s.adapters[t] = adapter
}
func (s *Registry) Adapter(t reflect.Type) Adapter {
	if s.adapters == nil {
		return nil
	}
	return s.adapters[t]
}
func (s *Registry) Read(t reflect.Type, r *Reader) any {
	if adapter := s.Adapter(t); adapter != nil {
		return adapter.Read(*s, *r)
	}
	panic("jso: no adapter registered for type " + t.String())
}
func (s *Registry) Write(w *Writer, v any) {
	t := reflect.TypeOf(v)
	if adapter := s.Adapter(t); adapter != nil {
		adapter.Write(*s, *w, v)
		return
	}
	panic("jso: no adapter registered for type " + t.String())
}
func (s *Registry) WriteWithType(w *Writer, v any, t reflect.Type) {
	if adapter := s.Adapter(t); adapter != nil {
		adapter.Write(*s, *w, v)
		return
	}
	panic("jso: no adapter registered for type " + t.String())
}
