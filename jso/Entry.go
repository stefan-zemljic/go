package jso

type Entry[K comparable, V any] struct {
	Key   K
	Value V
}
