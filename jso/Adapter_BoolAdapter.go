package jso

type BoolAdapter struct{}

var _ Adapter = BoolAdapter{}

func (a BoolAdapter) Write(_ *Registry, buffer *Buffer, value any) {
	if _, ok := value.(bool); ok {
		buffer.Add(value)
		return
	}
	panic("expected bool")
}
func (a BoolAdapter) Read(_ *Registry, data *Data) any {
	return data.MustBool()
}
