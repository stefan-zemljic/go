package jso

import (
	"reflect"
)

type Adapter interface {
	Read(reflect.Type, Stream) (any, bool)
	Write(any, Stream) bool
}
