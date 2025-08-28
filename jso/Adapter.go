package jso

type Adapter interface {
	Write(Registry, Writer, any)
	Read(Registry, Reader) any
}
