package shallow

import (
	"reflect"
	"testing"
	"unsafe"
)

type Src struct {
	A int
	B string
}

type Dst struct {
	A int
	B string
}

type DstExtra struct {
	A int
	B string
	C int
}

type DstMismatch struct {
	A int
	B int
}

func TestMap_SuccessAndCache(t *testing.T) {
	src := Src{A: 42, B: "hello"}
	var dst Dst
	Map(src, &dst)
	if dst.A != 42 || dst.B != "hello" {
		t.Errorf("expected {42 hello}, got %+v", dst)
	}
	src2 := Src{A: 99, B: "cached"}
	var dst2 Dst
	Map(src2, &dst2)
	if dst2.A != 99 || dst2.B != "cached" {
		t.Errorf("expected {99 cached}, got %+v", dst2)
	}
}

func TestMap_FromNotStructPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	var dst Dst
	Map(123, &dst)
}

func TestMap_ToNotPointerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	src := Src{}
	dst := Dst{}
	Map(src, dst)
}

func TestMap_ToNotStructPointerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	src := Src{}
	var dst *int
	Map(src, &dst)
}

func TestMap_FieldNotFoundPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	src := Src{}
	var dst DstExtra
	Map(src, &dst)
}

func TestMap_FieldTypeMismatchPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	src := Src{}
	var dst DstMismatch
	Map(src, &dst)
}

func TestMemmoveCopiesData(t *testing.T) {
	src := []byte{1, 2, 3, 4}
	dst := make([]byte, 4)
	memmove(
		unsafe.Pointer(&dst[0]),
		unsafe.Pointer(&src[0]),
		uintptr(len(src)),
	)
	if !reflect.DeepEqual(src, dst) {
		t.Errorf("memmove failed, expected %v, got %v", src, dst)
	}
}
