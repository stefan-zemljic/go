package openapi

import (
	"github.com/stefan-zemljic/go/jso"
)

type Builder struct {
	data jso.Map[string, any]
}

func (b *Builder) Build() []byte {
	json := b.data
}
