package z

import (
	"testing"
)

func TestBoxOf(t *testing.T) {
	b1 := BoxOf(42, true)
	if b1.Empty().V {
		t.Error("expected non-empty")
	}
	if v, ok := b1.Get(); !ok || v != 42 {
		t.Errorf("Get mismatch: %v %v", v, ok)
	}
	b2 := BoxOf(42, false)
	if !b2.Empty().V {
		t.Error("expected empty")
	}
	if v, ok := b2.Get(); ok || v != 0 {
		t.Errorf("Get mismatch: %v %v", v, ok)
	}
}

func TestBoxEmptyAndJust(t *testing.T) {
	e := Empty[int]()
	if !e.Empty().V {
		t.Error("expected empty")
	}
	if e.Present().V {
		t.Error("expected not present")
	}
	j := Just("hi")
	if j.Empty().V {
		t.Error("expected not empty")
	}
	if !j.Present().V {
		t.Error("expected present")
	}
}

func TestBoxMust(t *testing.T) {
	b := Just(7)
	if b.Must() != 7 {
		t.Error("Must returned wrong value")
	}
	mustPanic(t, func() { Empty[int]().Must() }, errorPrefixBox+"box empty")
}

func TestBoxMustNot(t *testing.T) {
	mustPanic(t, func() { Just(1).MustNot() }, errorPrefixBox+"box not empty")
	Empty[int]().MustNot()
}

func TestBoxMustNonZero(t *testing.T) {
	mustPanic(t, func() { Just(5).MustNonZero() }, errorPrefixBox+"box empty")
	mustPanic(t, func() { Just(0).MustNonZero() }, errorPrefixBox+"box zero")
	mustPanic(t, func() { Empty[int]().MustNonZero() }, errorPrefixBox+"box empty")
}

func TestBoxMustZero(t *testing.T) {
	Just(0).MustZero()
	mustPanic(t, func() { Just(3).MustZero() }, errorPrefixBox+"box not zero")
	Empty[int]().MustZero()
}

func TestBoxOrZero(t *testing.T) {
	if got := Just(9).OrZero(); got != 9 {
		t.Errorf("OrZero got %v", got)
	}
	if got := Empty[int]().OrZero(); got != 0 {
		t.Errorf("OrZero got %v", got)
	}
}

func TestBoxOrNil(t *testing.T) {
	if got := Just(10).OrNil(); got == nil || *got != 10 {
		t.Errorf("OrNil mismatch %v", got)
	}
	if got := Empty[int]().OrNil(); got != nil {
		t.Errorf("OrNil expected nil, got %v", got)
	}
}

func TestBoxOr(t *testing.T) {
	if got := Just("x").Or("y"); got != "x" {
		t.Errorf("Or got %v", got)
	}
	if got := Empty[string]().Or("y"); got != "y" {
		t.Errorf("Or got %v", got)
	}
}

func TestBoxReflectCoverage(t *testing.T) {
	type S struct{ A int }
	b := Just(S{})
	mustPanic(t, func() { b.MustNonZero() }, errorPrefixBox+"box zero")
}
