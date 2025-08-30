package z

import (
	"fmt"
)

type Res[T any] struct {
	V T
	E error
}

func Ok[T any](data T) Res[T] {
	return Res[T]{data, nil}
}

func Err[T any](err error) Res[T] {
	var zero T
	return Res[T]{zero, err}
}

func Errorf[T any](format string, args ...any) Res[T] {
	var zero T
	return Res[T]{zero, fmt.Errorf(format, args...)}
}

func ResOf[T any](data T, err error) Res[T] {
	return Res[T]{data, err}
}

func (r Res[T]) Ok() Flag {
	return FlagOf(r.E == nil)
}

func (r Res[T]) Err() Flag {
	return FlagOf(r.E != nil)
}

func (r Res[T]) Error() error {
	return r.E
}

func (r Res[T]) Get() (T, error) {
	return r.V, r.E
}

func (r Res[T]) Must() T {
	if r.E != nil {
		panic(fmt.Errorf("%s%w", errorPrefixRes, r.E))
	}
	return r.V
}

func (r Res[T]) MustErr() error {
	if r.E == nil {
		panic(errorPrefixRes + "no error")
	}
	return r.E
}

func (r Res[T]) Or(other T) T {
	if r.E != nil {
		return other
	}
	return r.V
}

func (r Res[T]) OrZero() T {
	if r.E != nil {
		var zero T
		return zero
	}
	return r.V
}

const errorPrefixRes = "z.Res: "
