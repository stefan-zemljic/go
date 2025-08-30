package z

import (
	"reflect"
)

type Box[T any] struct {
	V T
	P bool
}

func BoxOf[T any](data T, present bool) Box[T] {
	if present {
		return Box[T]{data, true}
	}
	var zero T
	return Box[T]{zero, false}
}

func Empty[T any]() Box[T] {
	var zero T
	return Box[T]{zero, false}
}

func Just[T any](data T) Box[T] {
	return Box[T]{data, true}
}

func (s Box[T]) Empty() Flag {
	return FlagOf(!s.P)
}

func (s Box[T]) Present() Flag {
	return FlagOf(s.P)
}

func (s Box[T]) Get() (T, bool) {
	return s.V, s.P
}

func (s Box[T]) Must() T {
	if !s.P {
		panic(errorPrefixBox + "box empty")
	}
	return s.V
}

func (s Box[T]) MustNot() {
	if s.P {
		panic(errorPrefixBox + "box not empty")
	}
}

func (s Box[T]) MustNonZero() T {
	if s.P {
		if reflect.ValueOf(s.V).IsZero() {
			panic(errorPrefixBox + "box zero")
		}
	}
	panic(errorPrefixBox + "box empty")
}

func (s Box[T]) MustZero() {
	if s.P {
		if !reflect.ValueOf(s.V).IsZero() {
			panic(errorPrefixBox + "box not zero")
		}
	}
}

func (s Box[T]) OrZero() T {
	if !s.P {
		var zero T
		return zero
	}
	return s.V
}

func (s Box[T]) OrNil() *T {
	if s.P {
		v := s.V
		return &v
	}
	return nil
}

func (s Box[T]) Or(other T) T {
	if s.P {
		return s.V
	}
	return other
}

const errorPrefixBox = "z.Box: "
