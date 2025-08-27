package jso

import (
	"slices"
)

type MultiMap[K comparable, V any] struct {
	entries []Entry[K, V]
	values  map[K][]V
}

func (s *MultiMap[K, V]) Keys() []K {
	return filterMap(s.entries, func(e Entry[K, V]) (K, bool) {
		return e.Key, true
	})
}
func (s *MultiMap[K, V]) Values() []V {
	return filterMap(s.entries, func(e Entry[K, V]) (V, bool) {
		return e.Value, true
	})
}
func (s *MultiMap[K, V]) Entries() []Entry[K, V] {
	return slices.Clone(s.entries)
}
func (s *MultiMap[K, V]) Map() map[K][]V {
	multiMap := map[K][]V{}
	for key, values := range s.values {
		multiMap[key] = slices.Clone(values)
	}
	return multiMap
}
func (s *MultiMap[K, V]) Get(key K) []V {
	if len(s.entries) > 0 {
		return slices.Clone(s.values[key])
	}
	return nil
}
func (s *MultiMap[K, V]) Put(key K, value V) {
	if s.values == nil {
		s.values = map[K][]V{key: {value}}
		s.entries = []Entry[K, V]{{Key: key, Value: value}}
	} else if values, exists := s.values[key]; exists {
		s.values[key] = append(values, value)
		s.entries = append(s.entries, Entry[K, V]{key, value})
	} else {
		s.values[key] = []V{value}
		s.entries = append(s.entries, Entry[K, V]{key, value})
	}
}
func (s *MultiMap[K, V]) Delete(key K) {
	if len(s.entries) > 0 {
		delete(s.values, key)
		s.entries = slices.DeleteFunc(s.entries, func(e Entry[K, V]) bool {
			return e.Key == key
		})
	}
}
func (s *MultiMap[K, V]) Clear() {
	s.entries = nil
	s.entries = nil
	s.values = nil
	s.values = nil
}
func (s *MultiMap[K, V]) ValueCount() int {
	return len(s.entries)
}
func (s *MultiMap[K, V]) KeyCount() int {
	return len(s.values)
}
func (s *MultiMap[K, V]) IsEmpty() bool {
	return len(s.entries) == 0
}
func (s *MultiMap[K, V]) Clone() *MultiMap[K, V] {
	values := map[K][]V{}
	for key, vs := range s.values {
		values[key] = slices.Clone(vs)
	}
	clone := &MultiMap[K, V]{
		entries: slices.Clone(s.entries),
		values:  values,
	}
	return clone
}
func (s *MultiMap[K, V]) Iter(fn func(K, V) bool) {
	for _, entry := range s.entries {
		if !fn(entry.Key, entry.Value) {
			break
		}
	}
}
