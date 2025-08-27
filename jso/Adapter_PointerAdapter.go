package jso

import (
	"fmt"
	"reflect"
)

type PointerAdapter[T any] struct {
	writeNil bool
	readNil  bool
}

var _ Adapter = PointerAdapter[int]{}

func (a PointerAdapter[T]) Write(registry *Registry, buffer *Buffer, value any) {
	v, ok := value.(*T)
	if !ok {
		panic(fmt.Sprintf("expected type *%T value, got type %T", *new(T), value))
	} else if v == nil {
		if a.writeNil {
			buffer.Add(nil)
			return
		}
		panic("nil pointer encountered, but writeNil is false")
	}
	registry.Write(reflect.TypeFor[T](), buffer, *v)
}
func (a PointerAdapter[T]) Read(registry *Registry, data *Data) any {
	if data.Null() {
		if a.readNil {
			return (*T)(nil)
		}
		panic("nil pointer encountered, but readNil is false")
	}
	v := registry.Read(reflect.TypeFor[T](), data)
	return &v
}
