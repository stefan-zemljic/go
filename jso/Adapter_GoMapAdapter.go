package jso

import (
	"fmt"
	"maps"
	"reflect"
)

type GoMapAdapter[K comparable, V any] struct {
	writeNil        bool
	writeEmptyAsNil bool
	writeNilAsEmpty bool
	readNil         bool
	readEmptyAsNil  bool
	readNilAsEmpty  bool
}

var _ Adapter = GoMapAdapter[int, int]{}

func (a GoMapAdapter[K, V]) Write(registry *Registry, buf *Buffer, value any) {
	v, ok := value.(map[K]V)
	if !ok {
		panic(fmt.Sprintf("expected type %T, got type %T", map[K]V{}, value))
	} else if v == nil {
		if a.writeNilAsEmpty {
			v = map[K]V{}
		} else if a.writeNil {
			buf.Nil()
			return
		}
		panic("nil map encountered, but writeNil is false")
	} else if len(v) == 0 && a.writeEmptyAsNil {
		if a.writeNil {
			buf.Nil()
			return
		}
		panic("writeEmptyAsNil, but writeNil is false")
	}
	buf.Obj()
	for key := range maps.Keys(v) {
		registry.Write(reflect.TypeFor[K](), buf, key)
		registry.Write(reflect.TypeFor[V](), buf, v[key])
	}
	buf.End()
}

func (a GoMapAdapter[K, V]) Read(registry *Registry, data *Data) any {
	if data.Null() {
		if a.readNilAsEmpty {
			return map[K]V{}
		} else if a.readNil {
			return nil
		}
		panic("nil map encountered, but readNil is false")
	}
	data.MustObject()
	if data.End() {
		if a.readEmptyAsNil {
			if a.readNil {
				return nil
			}
			panic("readEmptyAsNil, but readNil is false")
		}
		return map[K]V{}
	}
	m := map[K]V{}
	for {
		key := registry.Read(reflect.TypeFor[K](), data)
		value := registry.Read(reflect.TypeFor[V](), data)
		m[key.(K)] = value.(V)
		if data.End() {
			break
		}
	}
	return m
}
