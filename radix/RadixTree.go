package radix

import (
	rad "github.com/armon/go-radix"
)

type WalkFn[V any] func(s string, v V) bool
type Tree[V any] struct {
	tree rad.Tree
}

func New[V any]() *Tree[V] {
	return &Tree[V]{tree: *rad.New()}
}
func (s *Tree[V]) Len() int {
	return s.tree.Len()
}
func (s *Tree[V]) Insert(st string, v V) (V, bool) {
	if x, ok := s.tree.Insert(st, v); ok {
		return x.(V), true
	}
	var zero V
	return zero, false
}
func (s *Tree[V]) Delete(st string) (V, bool) {
	if x, ok := s.tree.Delete(st); ok {
		return x.(V), true
	}
	var zero V
	return zero, false
}
func (s *Tree[V]) DeletePrefix(st string) int {
	return s.tree.DeletePrefix(st)
}
func (s *Tree[V]) Get(st string) (V, bool) {
	if x, ok := s.tree.Get(st); ok {
		return x.(V), true
	}
	var zero V
	return zero, false
}
func (s *Tree[V]) LongestPrefix(st string) (string, V, bool) {
	if str, x, ok := s.tree.LongestPrefix(st); ok {
		return str, x.(V), true
	}
	var zero V
	return "", zero, false
}
func (s *Tree[V]) Minimum() (string, V, bool) {
	if str, x, ok := s.tree.Minimum(); ok {
		return str, x.(V), true
	}
	var zero V
	return "", zero, false
}
func (s *Tree[V]) Maximum() (string, V, bool) {
	if str, x, ok := s.tree.Maximum(); ok {
		return str, x.(V), true
	}
	var zero V
	return "", zero, false
}
func (s *Tree[V]) Walk(fn WalkFn[V]) {
	s.tree.Walk(func(st string, v any) bool {
		return !fn(st, v.(V))
	})
}
func (s *Tree[V]) WalkPrefix(prefix string, fn WalkFn[V]) {
	s.tree.WalkPrefix(prefix, func(st string, v any) bool {
		return !fn(st, v.(V))
	})
}
func (s *Tree[V]) WalkPath(path string, fn WalkFn[V]) {
	s.tree.WalkPath(path, func(st string, v any) bool {
		return !fn(st, v.(V))
	})
}
func (s *Tree[V]) ToMap() map[string]V {
	m := map[string]V{}
	s.tree.Walk(func(st string, v any) bool {
		m[st] = v.(V)
		return false
	})
	return m
}
