package jso

type Adapter interface {
	Write(registry *Registry, buffer *Buffer, value any)
	Read(registry *Registry, data *Data) any
}
