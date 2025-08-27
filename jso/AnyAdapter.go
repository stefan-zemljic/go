package jso

import (
	"reflect"
)

type AnyAdapter struct{}

func (a AnyAdapter) Write(registry *Registry, buf *Buffer, value any) {
	if value == nil {
		buf.Nil()
		return
	}
	name := reflect.TypeOf(value).String()
	registry.Type(name)
	buf.Arr()
	buf.Add(name)
	registry.Write(reflect.TypeFor[any](), buf, value)
	buf.End()
}
func (a AnyAdapter) Read(registry *Registry, data *Data) any {
	if data.Null() {
		return nil
	}
	data.MustArray()
	typeName := data.MustString()
	value := registry.Read(registry.Type(typeName), data)
	data.MustEnd()
	return value
}
