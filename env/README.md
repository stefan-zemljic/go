# env

`env` is a tiny Go library for safely generating shell scripts that set or unset environment variables.
Instead of directly interpolating environment values (which can be unsafe), it uses a `hex2str` helper to encode values in hex, ensuring proper escaping and safe export/unset operations.

---

## Features

* Generate shell code to **set environment variables** with safe hex escaping.
* Generate shell code to **unset environment variables**.
* Collects all output in a buffer, retrievable as bytes, string, or written to an `io.Writer`.
* Simple API with reset and write support.

---

## Installation

```bash
github.com/stefan-zemljic/go/env
```

---

## Usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/stefan-zemljic/go/env"
)

func main() {
	var e env.Emitter

	// Set an environment variable
	e.Set("GREETING", "hello world")

	// Unset another
	e.Unset("OLD_VAR")

	// Print generated shell script
	fmt.Println(e.String())

	// Or write directly to stdout
	e.WriteTo(os.Stdout)
}
```

The above will produce something like:

```sh
hex2str() {
  local hex="$1"
  local esc=""
  local i
  for ((i=0; i<${#hex}; i+=2)); do
    esc+="\\x${hex:$i:2}"
  done
  eval "echo \\$'$esc'"
}
export GREETING=$(hex2str 68656c6c6f20776f726c64)
unset OLD_VAR
```

---

## API

### `Emitter`

* `func (e *Emitter) Set(name, value string)`
  Adds a command to set an environment variable with safe encoding.

* `func (e *Emitter) Unset(name string)`
  Adds a command to unset an environment variable.

* `func (e *Emitter) Bytes() []byte`
  Returns the accumulated script as a byte slice.

* `func (e *Emitter) String() string`
  Returns the accumulated script as a string.

* `func (e *Emitter) WriteTo(w io.Writer) (int64, error)`
  Writes the script to an `io.Writer`.

* `func (e *Emitter) Reset()`
  Clears the buffer.

* `func (e *Emitter) Write(bs []byte) (int, error)`
  Appends raw bytes to the buffer.

---

## License

MIT License

```
Copyright (c) 2025 Stefan Zemljic

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```
