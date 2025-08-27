package jso

type Float interface {
	~float32 | ~float64
}
type FloatAdapter[T Float] struct {
	bitSize int
}

var _ Adapter = FloatAdapter[float32]{}

func (a FloatAdapter[T]) Write(_ *Registry, buffer *Buffer, value any) {
	if _, ok := value.(T); ok {
		buffer.Add(value)
	} else {
		panic("expected float")
	}
}
func (a FloatAdapter[T]) Read(_ *Registry, data *Data) any {
	return mustFloat[T](data, a.bitSize)
}
