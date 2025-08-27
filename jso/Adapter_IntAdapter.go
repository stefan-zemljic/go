package jso

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type IntAdapter[T Int] struct {
	bitSize int
}

var _ Adapter = IntAdapter[int]{}

func (a IntAdapter[T]) Write(_ *Registry, buffer *Buffer, value any) {
	if _, ok := value.(T); ok {
		buffer.Add(value)
	} else {
		panic("expected int")
	}
}
func (a IntAdapter[T]) Read(_ *Registry, data *Data) any {
	return mustInt[T](data, a.bitSize)
}
