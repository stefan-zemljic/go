package jso

import (
	"reflect"
)

type Registry struct {
	types    map[string]reflect.Type
	adapters map[reflect.Type]Adapter
}

func New() *Registry {
	r := &Registry{
		adapters: map[reflect.Type]Adapter{
			reflect.TypeFor[any]():                    AnyAdapter{},
			reflect.TypeFor[bool]():                   BoolAdapter{},
			reflect.TypeFor[int]():                    IntAdapter[int]{0},
			reflect.TypeFor[int8]():                   IntAdapter[int8]{8},
			reflect.TypeFor[int16]():                  IntAdapter[int16]{16},
			reflect.TypeFor[int32]():                  IntAdapter[int32]{32},
			reflect.TypeFor[int64]():                  IntAdapter[int64]{64},
			reflect.TypeFor[uint]():                   UintAdapter[uint]{0},
			reflect.TypeFor[uint8]():                  UintAdapter[uint8]{8},
			reflect.TypeFor[uint16]():                 UintAdapter[uint16]{16},
			reflect.TypeFor[uint32]():                 UintAdapter[uint32]{32},
			reflect.TypeFor[uint64]():                 UintAdapter[uint64]{64},
			reflect.TypeFor[uintptr]():                UintAdapter[uintptr]{0},
			reflect.TypeFor[float32]():                FloatAdapter[float32]{64},
			reflect.TypeFor[float64]():                FloatAdapter[float64]{64},
			reflect.TypeFor[string]():                 StringAdapter{},
			reflect.TypeFor[map[any]any]():            GoMapAdapter[any, any]{writeNil: true, readNil: true},
			reflect.TypeFor[Map[any, any]]():          MapAdapter[any, any]{writeNil: true, readNil: true},
			reflect.TypeFor[MultiMap[any, any]]():     MultiMapAdapter[any, any]{writeNil: true, readNil: true},
			reflect.TypeFor[CombinedMap[any, any]]():  CombinedMapAdapter[any, any]{writeNil: true, readNil: true},
			reflect.TypeFor[[]any]():                  SliceAdapter[any]{writeNil: true, readNil: true},
			reflect.TypeFor[*any]():                   PointerAdapter[any]{writeNil: true, readNil: true},
			reflect.TypeFor[*bool]():                  PointerAdapter[bool]{writeNil: true, readNil: true},
			reflect.TypeFor[*int]():                   PointerAdapter[int]{writeNil: true, readNil: true},
			reflect.TypeFor[*int8]():                  PointerAdapter[int8]{writeNil: true, readNil: true},
			reflect.TypeFor[*int16]():                 PointerAdapter[int16]{writeNil: true, readNil: true},
			reflect.TypeFor[*int32]():                 PointerAdapter[int32]{writeNil: true, readNil: true},
			reflect.TypeFor[*int64]():                 PointerAdapter[int64]{writeNil: true, readNil: true},
			reflect.TypeFor[*uint]():                  PointerAdapter[uint]{writeNil: true, readNil: true},
			reflect.TypeFor[*uint8]():                 PointerAdapter[uint8]{writeNil: true, readNil: true},
			reflect.TypeFor[*uint16]():                PointerAdapter[uint16]{writeNil: true, readNil: true},
			reflect.TypeFor[*uint32]():                PointerAdapter[uint32]{writeNil: true, readNil: true},
			reflect.TypeFor[*uint64]():                PointerAdapter[uint64]{writeNil: true, readNil: true},
			reflect.TypeFor[*uintptr]():               PointerAdapter[uintptr]{writeNil: true, readNil: true},
			reflect.TypeFor[*float32]():               PointerAdapter[float32]{writeNil: true, readNil: true},
			reflect.TypeFor[*float64]():               PointerAdapter[float64]{writeNil: true, readNil: true},
			reflect.TypeFor[*string]():                PointerAdapter[string]{writeNil: true, readNil: true},
			reflect.TypeFor[*map[any]any]():           PointerAdapter[map[any]any]{writeNil: true, readNil: true},
			reflect.TypeFor[*Map[any, any]]():         PointerAdapter[Map[any, any]]{writeNil: true, readNil: true},
			reflect.TypeFor[*MultiMap[any, any]]():    PointerAdapter[MultiMap[any, any]]{writeNil: true, readNil: true},
			reflect.TypeFor[*CombinedMap[any, any]](): PointerAdapter[CombinedMap[any, any]]{writeNil: true, readNil: true},
		},
	}
	return r
}
func (s *Registry) Adapter(t reflect.Type, adapter Adapter, overwrite ...bool) {
	if _, ok := s.adapters[t]; ok && !append(overwrite, false)[0] {
		panic("Adapter already registered for type " + t.String())
	}
	s.types[t.String()] = t
	s.adapters[t] = adapter
}
func (s *Registry) Write(t reflect.Type, buffer *Buffer, value any) {
	if writer, ok := s.adapters[t]; ok {
		writer.Write(s, buffer, value)
		return
	}
	panic("no adapter registered for type " + t.String())
}
func (s *Registry) Type(name string) reflect.Type {
	if t, ok := s.types[name]; ok {
		return t
	}
	panic("no type registered for name " + name)
}
func (s *Registry) Read(t reflect.Type, data *Data) any {
	if reader, ok := s.adapters[t]; ok {
		return reader.Read(s, data)
	}
	panic("no adapter registered for type " + t.String())
}
func (s *Registry) AdapterFor(t reflect.Type) Adapter {
	if adapter, ok := s.adapters[t]; ok {
		return adapter
	}
	return nil
}
