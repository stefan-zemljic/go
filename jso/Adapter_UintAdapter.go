package jso

type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
type UintAdapter[T Uint] struct {
	bitSize int
}

var _ Adapter = UintAdapter[uint]{}

func (a UintAdapter[T]) Write(_ *Registry, buffer *Buffer, value any) {
	if v, ok := value.(T); ok {
		buffer.Add(v)
		return
	}
	panic("expected uint")
}
func (a UintAdapter[T]) Read(_ *Registry, data *Data) any {
	return mustUint[T](data, a.bitSize)
}
