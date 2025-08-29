package jso

import (
	"fmt"
	"testing"
)

func wrapInt[T any](fn func(*Reader) (T, bool)) func(*Reader) (string, bool) {
	return func(r *Reader) (string, bool) {
		v, ok := fn(r)
		return fmt.Sprint(v), ok
	}
}

func TestReader_NonMustVariants(t *testing.T) {
	tests := []struct {
		name string
		json string
		fn   any
		want string
	}{
		{"Int", "123", wrapInt((*Reader).Int), "123"},
		{"Int8", "8", wrapInt((*Reader).Int8), "8"},
		{"Int16", "16", wrapInt((*Reader).Int16), "16"},
		{"Int32", "32", wrapInt((*Reader).Int32), "32"},
		{"Int64", "64", wrapInt((*Reader).Int64), "64"},
		{"Uint", "7", wrapInt((*Reader).Uint), "7"},
		{"Uint8", "8", wrapInt((*Reader).Uint8), "8"},
		{"Uint16", "16", wrapInt((*Reader).Uint16), "16"},
		{"Uint32", "32", wrapInt((*Reader).Uint32), "32"},
		{"Uint64", "64", wrapInt((*Reader).Uint64), "64"},
		{"Uintptr", "99", wrapInt((*Reader).Uintptr), "99"},
		{"Float32", "1.25", wrapInt((*Reader).Float32), "1.25"},
		{"Float64", "2.5", wrapInt((*Reader).Float64), "2.5"},
		{"Number", "123", wrapInt((*Reader).Number), "123"},
		{"String", `"xyz"`, wrapInt((*Reader).String), "xyz"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reader{buf: []byte(tt.json)}
			var got string
			var ok bool
			switch fn := tt.fn.(type) {
			case func(*Reader) (string, bool):
				got, ok = fn(r)
			default:
				t.Fatalf("unsupported function type %T", tt.fn)
			}
			if !ok {
				t.Fatalf("expected ok")
			}
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReader_Success(t *testing.T) {
	tests := []struct {
		name string
		json string
		fn   func(*Reader) string
		want string
	}{
		{"null", "null", func(r *Reader) string {
			if !r.Null() {
				t.Fatal("expected null")
			}
			return "null"
		}, "null"},
		{"bool true", "true", func(r *Reader) string {
			if v, ok := r.Bool(); !ok || !v {
				t.Fatal("expected true")
			}
			return "true"
		}, "true"},
		{"bool false", "false", func(r *Reader) string {
			if r.MustBool() {
				t.Fatal("expected false")
			}
			return "false"
		}, "false"},
		{"int", "123", func(r *Reader) string {
			return fmt.Sprint(r.MustInt())
		}, "123"},
		{"int8", "8", func(r *Reader) string {
			return fmt.Sprint(r.MustInt8())
		}, "8"},
		{"int16", "16", func(r *Reader) string {
			return fmt.Sprint(r.MustInt16())
		}, "16"},
		{"int32", "32", func(r *Reader) string {
			return fmt.Sprint(r.MustInt32())
		}, "32"},
		{"int64", "64", func(r *Reader) string {
			return fmt.Sprint(r.MustInt64())
		}, "64"},

		{"uint", "7", func(r *Reader) string {
			return fmt.Sprint(r.MustUint())
		}, "7"},
		{"uint8", "8", func(r *Reader) string {
			return fmt.Sprint(r.MustUint8())
		}, "8"},
		{"uint16", "16", func(r *Reader) string {
			return fmt.Sprint(r.MustUint16())
		}, "16"},
		{"uint32", "32", func(r *Reader) string {
			return fmt.Sprint(r.MustUint32())
		}, "32"},
		{"uint64", "64", func(r *Reader) string {
			return fmt.Sprint(r.MustUint64())
		}, "64"},
		{"uintptr", "99", func(r *Reader) string {
			return fmt.Sprint(r.MustUintptr())
		}, "99"},
		{"float32", "3.5", func(r *Reader) string {
			return fmt.Sprint(r.MustFloat32())
		}, "3.5"},
		{"float64", "12.5", func(r *Reader) string {
			return fmt.Sprint(r.MustFloat64())
		}, "12.5"},
		{"number as string", "42", func(r *Reader) string {
			return r.MustNumber()
		}, "42"},
		{"string", `"abc"`, func(r *Reader) string {
			return r.MustString()
		}, "abc"},
		{"array", `[1,2]`, func(r *Reader) string {
			if !r.Array() {
				t.Fatal("expected array")
			}
			r.MustInt()
			r.MustInt()
			if !r.End() {
				t.Fatal("expected end")
			}
			return "ok"
		}, "ok"},
		{"object", `{"a":1}`, func(r *Reader) string {
			if !r.Object() {
				t.Fatal("expected object")
			} else if k := r.MustString(); k != "a" {
				t.Fatalf("got %s", k)
			} else if v := r.MustInt(); v != 1 {
				t.Fatal("expected 1")
			}
			r.MustEnd()
			return "ok"
		}, "ok"},
		{"Save/Restore", `true false`, func(r *Reader) string {
			save := r.Save()
			if !r.MustBool() {
				t.Fatal("expected true")
			}
			r.MustBool()
			r.Restore(save)
			if !r.MustBool() {
				t.Fatal("expected restored true")
			}
			return "ok"
		}, "ok"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reader{buf: []byte(tt.json)}
			if got := tt.fn(r); got != tt.want {
				t.Errorf("%s got %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestReader_Errors(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		fn      func(r *Reader)
		wantErr string
	}{
		{"MustNull wrong", "true", func(r *Reader) { r.MustNull() }, "jso: expected null"},
		{"MustNull EOF", "", func(r *Reader) { r.MustNull() }, "jso: unexpected end of input"},
		{"MustBool wrong", "123", func(r *Reader) { r.MustBool() }, "jso: expected different token"},
		{"MustInt invalid", `"abc"`, func(r *Reader) { r.MustInt() }, "jso: expected different token"},
		{"MustFloat64 invalid num", `"abc"`, func(r *Reader) { r.MustFloat64() }, "jso: expected different token"},
		{"MustFloat64 parse error", `"12x"`, func(r *Reader) { r.next = num([]byte("12x")); r.has = true; r.MustFloat64() }, "jso: invalid number"},
		{"MustNumber wrong", "true", func(r *Reader) { r.MustNumber() }, "jso: expected different token"},
		{"MustString wrong", "123", func(r *Reader) { r.MustString() }, "jso: expected different token"},
		{"MustArray wrong", "123", func(r *Reader) { r.MustArray() }, "jso: expected different token"},
		{"MustObject wrong", "123", func(r *Reader) { r.MustObject() }, "jso: expected different token"},
		{"MustEnd wrong", "123", func(r *Reader) { r.MustEnd() }, "jso: expected end of array or object"},
		{"prepareNext invalid short null", "n", func(r *Reader) { r.prepareNext() }, "jso: invalid JSON input"},
		{"prepareNext invalid short true", "t", func(r *Reader) { r.prepareNext() }, "jso: invalid JSON input"},
		{"prepareNext invalid short false", "f", func(r *Reader) { r.prepareNext() }, "jso: invalid JSON input"},
		{"prepareNext invalid num", "x", func(r *Reader) { r.prepareNext() }, "jso: invalid JSON input"},
		{"prepareNext invalid string", `"bad`, func(r *Reader) { r.prepareNext() }, "jso: invalid JSON input"},
		{"Kind invalid state", "", func(r *Reader) {
			r.has = true
			r.next = struct{}{}
			_ = r.Kind()
		}, "jso: invalid state"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if rcv := recover(); rcv == nil {
					t.Fatalf("expected panic %q, got none", tt.wantErr)
				} else if got := fmt.Sprint(rcv); got != tt.wantErr {
					t.Fatalf("expected %q, got %q", tt.wantErr, got)
				}
			}()
			r := &Reader{buf: []byte(tt.json)}
			tt.fn(r)
		})
	}
}
