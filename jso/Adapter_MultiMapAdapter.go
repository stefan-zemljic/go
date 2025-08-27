package jso

import (
	"fmt"
	"reflect"
)

type MultiMapAdapter[K comparable, V any] struct {
	writeNil        bool
	writeEmptyAsNil bool
	writeNilAsEmpty bool
	readNil         bool
	readEmptyAsNil  bool
	readNilAsEmpty  bool
	writeLatest     bool
	writeFlat       bool
}

var _ Adapter = MultiMapAdapter[int, int]{}

func (a MultiMapAdapter[K, V]) Write(registry *Registry, buf *Buffer, value any) {
	v, ok := value.(*MultiMap[K, V])
	if !ok {
		panic(fmt.Sprintf("expected type *%T, got type %T", MultiMap[K, V]{}, value))
	} else if v == nil {
		if a.writeNilAsEmpty {
			v = &MultiMap[K, V]{}
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
	var entries []Entry[K, V]
	if a.writeLatest {
		entries = v.entries
	} else if a.writeFlat {
		entries = v.entries
	} else {
		var seen = map[K]struct{}{}
		for _, entry := range v.entries {
			if _, exists := seen[entry.Key]; !exists {
				seen[entry.Key] = struct{}{}
				registry.Write(reflect.TypeFor[K](), buf, entry.Key)
				buf.Arr()
				for _, val := range v.values[entry.Key] {
					registry.Write(reflect.TypeFor[V](), buf, val)
				}
				buf.End()
			}
		}
		buf.End()
		return
	}
	for _, entry := range entries {
		registry.Write(reflect.TypeFor[K](), buf, entry.Key)
		registry.Write(reflect.TypeFor[V](), buf, entry.Value)
	}
	buf.End()
}

func (a MultiMapAdapter[K, V]) Read(registry *Registry, data *Data) any {
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
	m := &MultiMap[K, V]{}
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
