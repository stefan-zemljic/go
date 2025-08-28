package shallow

import (
	"testing"
)

func TestReflectMap_SuccessAndCache(t *testing.T) {
	src := Src{A: 42, B: "hello"}
	var dst Dst
	reflectMap(src, &dst)
	if dst.A != 42 || dst.B != "hello" {
		t.Errorf("expected {42 hello}, got %+v", dst)
	}
	src2 := Src{A: 99, B: "cached"}
	var dst2 Dst
	reflectMap(src2, &dst2)
	if dst2.A != 99 || dst2.B != "cached" {
		t.Errorf("expected {99 cached}, got %+v", dst2)
	}
}

func TestReflectMap_FromNotStructPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	var dst Dst
	reflectMap(123, &dst)
}

func TestReflectMap_ToNotPointerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	src := Src{}
	dst := Dst{}
	reflectMap(src, dst)
}

func TestReflectMap_ToNotStructPointerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	src := Src{}
	var dst *int
	reflectMap(src, &dst)
}

func TestReflectMap_FieldNotFoundPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	src := Src{}
	var dst DstExtra
	reflectMap(src, &dst)
}

func TestReflectMap_FieldTypeMismatchPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	src := Src{}
	var dst DstMismatch
	reflectMap(src, &dst)
}
