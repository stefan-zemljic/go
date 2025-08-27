package jso

import (
	"fmt"
	"reflect"
)

type MapAdapter[K comparable, V any] struct {
	writeNil        bool
	writeEmptyAsNil bool
	writeNilAsEmpty bool
	readNil         bool
	readEmptyAsNil  bool
	readNilAsEmpty  bool
}

var _ Adapter = MapAdapter[int, int]{}

func (a MapAdapter[K, V]) Write(registry *Registry, buf *Buffer, value any) {
	v, ok := value.(*Map[K, V])
	if !ok {
		panic(fmt.Sprintf("expected type *%T, got type %T", Map[K, V]{}, value))
	} else if v == nil {
		if a.writeNilAsEmpty {
			v = &Map[K, V]{}
		} else if a.writeNil {
			buf.Nil()
			return
		}
		panic("nil map encountered, but writeNil is false")
	} else if len(v.entries) == 0 && a.writeEmptyAsNil {
		if a.writeNil {
			buf.Nil()
			return
		}
		panic("writeEmptyAsNil, but writeNil is false")
	}
	buf.Obj()
	for _, entry := range v.entries {
		registry.Write(reflect.TypeFor[K](), buf, entry.Key)
		registry.Write(reflect.TypeFor[V](), buf, entry.Value)
	}
	buf.End()
}

func (a MapAdapter[K, V]) Read(registry *Registry, data *Data) any {
	if data.Null() {
		if a.readNilAsEmpty {
			return &Map[K, V]{}
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
		return &Map[K, V]{}
	}
	m := &Map[K, V]{}
	for {
		key := registry.Read(reflect.TypeFor[K](), data)
		value := registry.Read(reflect.TypeFor[V](), data)
		m.Put(key.(K), value.(V))
		if data.End() {
			break
		}
	}
	return m
}
