package jso

type StringAdapter struct{}

var _ Adapter = StringAdapter{}

func (a StringAdapter) Write(_ *Registry, buffer *Buffer, value any) {
	if v, ok := value.(string); ok {
		buffer.Add(v)
		return
	}
	panic("expected string")
}
func (a StringAdapter) Read(_ *Registry, data *Data) any {
	return data.MustString()
}
