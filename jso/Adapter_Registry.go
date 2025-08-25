package jso

import (
	"fmt"
	"reflect"
)

type Registry struct {
	Adapters  MultiMap[reflect.Type, Adapter]
	Fallbacks []Adapter
}

func (r *Registry) Adapter(t reflect.Type, a Adapter) {
	r.Adapters.Put(t, a)
}
func (r *Registry) Fallback(a Adapter) {
	r.Fallbacks = append(r.Fallbacks, a)
}
func (r *Registry) Write(v any, s Stream) {
	t := reflect.TypeOf(v)
	for _, adapter := range r.Adapters.Get(t) {
		if adapter.Write(v, s) {
			return
		}
	}
	for _, adapter := range r.Fallbacks {
		if adapter.Write(v, s) {
			return
		}
	}
	panic(fmt.Sprintf("no adapter found for type %s with value %v", t.String(), v))
}
func (r *Registry) Read(t reflect.Type, s Stream) any {
	for _, adapter := range r.Adapters.Get(t) {
		if v, ok := adapter.Read(t, s); ok {
			return v
		}
	}
	for _, adapter := range r.Fallbacks {
		if v, ok := adapter.Read(t, s); ok {
			return v
		}
	}
	panic(fmt.Sprintf("no adapter found for type %s with payload %v", t.String(), s))
}
