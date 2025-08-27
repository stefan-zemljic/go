package jso

import (
	"fmt"
	"reflect"
)

type SliceAdapter[T any] struct {
	writeNil        bool
	writeNilAsEmpty bool
	writeEmptyAsNil bool
	readNil         bool
	readNilAsEmpty  bool
	readEmptyAsNil  bool
}

var _ Adapter = SliceAdapter[int]{}

func (a SliceAdapter[T]) Write(registry *Registry, buffer *Buffer, value any) {
	v, ok := value.([]T)
	if !ok {
		panic(fmt.Sprintf("expected type %T, got type %T", []T{}, value))
	} else if v == nil {
		if a.writeNilAsEmpty {
			v = []T{}
		} else if a.writeNil {
			buffer.Add(nil)
			return
		}
		panic("nil slice encountered, but writeNil is false")
	} else if len(v) == 0 && a.writeEmptyAsNil {
		if a.writeNil {
			buffer.Add(nil)
			return
		}
		panic("writeEmptyAsNil, but writeNil is false")
	}
	buffer.Arr()
	for _, item := range v {
		registry.Write(reflect.TypeFor[T](), buffer, item)
	}
	buffer.End()
}
func (a SliceAdapter[T]) Read(registry *Registry, data *Data) any {
	if data.Null() {
		if a.readNilAsEmpty {
			return []T{}
		} else if a.readNil {
			return ([]T)(nil)
		}
		panic("nil slice encountered, but readNil is false")
	}
	data.MustArray()
	if data.End() {
		if a.readEmptyAsNil {
			if a.readNil {
				return ([]T)(nil)
			}
			panic("readEmptyAsNil, but readNil is false")
		}
		return []T{}
	}
	slice := make([]T, 0)
	for !data.End() {
		item := registry.Read(reflect.TypeFor[T](), data)
		slice = append(slice, item.(T))
	}
	return slice
}
