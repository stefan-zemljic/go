# shallow

`shallow` provides fast struct-to-struct mapping in Go when the field names and types of the destination struct are a subset of the source struct.

It is useful when you need to copy data between types without writing repetitive manual code and want better performance than pure reflection.

Implemented for go 1.24, will likely break in future versions of Go.

No warranty is provided. Use at your own risk.

---

## Installation

```bash
go get github.com/stefan-zemljic/go/shallow
```

---

## Usage

Suppose you have two structs with matching fields:

```go
package main

import (
	"fmt"

	"github.com/stefan-zemljic/go/shallow"
)

type Src struct {
	A int
	B string
}

type Dst struct {
	A int
	B string
}

func main() {
	src := Src{A: 42, B: "hello"}
	var dst Dst

	shallow.Map(src, &dst)

	fmt.Printf("%+v\n", dst) // {A:42 B:hello}
}
```

### Rules

* `from` must be a **struct** or `*struct`.
* `to` must be a **pointer to a struct**.
* All fields in `to` must:

  * exist in `from` (matched by name), and
  * have the exact same type.

If these conditions are not met, `Map` panics with a descriptive error.

---

## Performance

This library implements two mapping strategies:

* **`Map` (default)** – uses `unsafe` + cached offsets for maximum speed.
* **`reflectMap` (internal)** – fallback implementation using `reflect` (mainly for comparison/benchmarks).

### Benchmarks

The following benchmarks were run with the dummy struct:

```go
type Foo struct {
	A int
	B string
	C float64
	D bool
}
```

Run:

```bash
go test -bench=.
```

On **Windows / amd64 / Go 1.24.6 / Intel i7-14700F**, the results were:

| Benchmark                 | Iterations | ns/op |
| ------------------------- | ---------- | ----- |
| Benchmark\_Map/Manual-28  | 548711145  | 2.172 |
| Benchmark\_Map/Unsafe-28  | 31381162   | 39.13 |
| Benchmark\_Map/Reflect-28 | 1621779    | 688.4 |

* **Manual mapping** is the fastest (essentially free).
* **Unsafe mapping** is more then 10x slower than manual, but requires no boilerplate code.
* **Reflect mapping** is significantly slower, and has no essential advantage over unsafe mapping.

---

## Testing

Run unit tests:

```bash
go test ./...
```

---

## Caveats

* Uses reflection **once** for each pair of from/to types to check types and cache mapping directives.
*  - Cache is `sync.Map`, safe for concurrent use, optimized for read-heavy workloads, matching this use case.
* **Shallow** copy of first level of values:
*  - Copies pointers as-is, does not deep-copy pointed-to values.
*  - `strings.Builder` fields are copied as-is, so if already initialized, they break.
* Casts any with `unsafe` to the current structure of `eface`, which is not guaranteed to be stable across Go versions.
* Is quite pedantic, types must match exactly, and no field of `to` can be missing in `from`.

---