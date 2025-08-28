package shallow

import (
	"testing"
)

type Foo struct {
	A int
	B string
	C float64
	D bool
}

func ManualMap(src Foo, dst *Foo) {
	dst.A = src.A
	dst.B = src.B
	dst.C = src.C
	dst.D = src.D
}

func Benchmark_Map(b *testing.B) {
	src := Foo{A: 42, B: "hello", C: 3.14, D: true}
	dst := Foo{}
	b.Run("Manual", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ManualMap(src, &dst)
		}
	})
	b.Run("Unsafe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			Map(src, &dst)
		}
	})
	b.Run("Reflect", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			reflectMap(src, &dst)
		}
	})
}
