package jso

import (
	"maps"
	"slices"
)

type Map[K comparable, V any] struct {
	entries []Entry[K, V]
	values  map[K]V
}

func (s *Map[K, V]) Keys() []K {
	return filterMap(s.entries, func(e Entry[K, V]) (K, bool) {
		return e.Key, true
	})
}
func (s *Map[K, V]) Values() []V {
	return filterMap(s.entries, func(e Entry[K, V]) (V, bool) {
		return e.Value, true
	})
}
func (s *Map[K, V]) Entries() []Entry[K, V] {
	return slices.Clone(s.entries)
}
func (s *Map[K, V]) Map() map[K]V {
	return maps.Clone(s.values)
}
func (s *Map[K, V]) Get(key K) (value V, ok bool) {
	if len(s.entries) > 0 {
		value, ok = s.values[key]
	}
	return
}
func (s *Map[K, V]) MustGet(key K) V {
	if value, ok := s.Get(key); ok {
		return value
	}
	var zero V
	return zero
}
func (s *Map[K, V]) Put(key K, value V) {
	if s.values == nil {
		s.values = map[K]V{key: value}
		s.entries = []Entry[K, V]{{Key: key, Value: value}}
	} else if _, exists := s.values[key]; exists {
		s.values[key] = value
		s.entries = slices.DeleteFunc(s.entries, func(e Entry[K, V]) bool {
			return e.Key == key
		})
		s.entries = append(s.entries, Entry[K, V]{key, value})
	} else {
		s.values[key] = value
		s.entries = append(s.entries, Entry[K, V]{key, value})
	}
}
func (s *Map[K, V]) Delete(key K) {
	if len(s.entries) > 0 {
		delete(s.values, key)
		s.entries = slices.DeleteFunc(s.entries, func(e Entry[K, V]) bool {
			return e.Key == key
		})
	}
}
func (s *Map[K, V]) Clear() {
	s.entries = nil
	s.values = nil
}
func (s *Map[K, V]) Len() int {
	return len(s.entries)
}
func (s *Map[K, V]) IsEmpty() bool {
	return len(s.entries) == 0
}
func (s *Map[K, V]) Clone() *Map[K, V] {
	clone := &Map[K, V]{
		entries: slices.Clone(s.entries),
		values:  maps.Clone(s.values),
	}
	return clone
}
func (s *Map[K, V]) Iter(fn func(K, V) bool) {
	for _, entry := range s.entries {
		if !fn(entry.Key, entry.Value) {
			break
		}
	}
}
func (s *Map[K, V]) ToWriter(w *Writer) {
	w.Add(Object)
	for _, entry := range s.entries {
		w.Add(entry.Key)
		w.Add(entry.Value)
	}
	w.Add(End)
}
func (s *Map[K, V]) To(st *Stream) {
	st.Add(Object)
	for _, entry := range s.entries {
		st.AddAll(entry.Key, entry.Value)
	}
	st.Add(End)
}
func (s *Map[K, V]) Pretty() string {
	var w Writer
	s.ToWriter(&w)
	return w.Pretty()
}
func (s *Map[K, V]) String() string {
	var w Writer
	s.ToWriter(&w)
	return w.String()
}
