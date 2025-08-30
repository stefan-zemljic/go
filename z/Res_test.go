package z

import (
	"errors"
	"fmt"
	"testing"
)

func mustPanicRes(t *testing.T, fn func(), want string) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			got := fmt.Sprint(r)
			if got != want {
				t.Fatalf("panic mismatch: got %q, want %q", got, want)
			}
		} else {
			t.Fatalf("expected panic %q, got none", want)
		}
	}()
	fn()
}

func TestResOkAndErrConstructors(t *testing.T) {
	r1 := Ok(123)
	if r1.E != nil || r1.V != 123 {
		t.Errorf("Ok wrong: %+v", r1)
	}
	r2 := Err[int](errors.New("fail"))
	if r2.E == nil || r2.V != 0 {
		t.Errorf("Err wrong: %+v", r2)
	}
	r3 := Errorf[int]("bad %d", 7)
	if r3.E == nil || r3.V != 0 || r3.E.Error() != "bad 7" {
		t.Errorf("Errorf wrong: %+v", r3)
	}
	r4 := ResOf("hi", nil)
	if r4.V != "hi" || r4.E != nil {
		t.Errorf("ResOf wrong: %+v", r4)
	}
}

func TestResOkErrFlags(t *testing.T) {
	if !Ok("x").Ok().V {
		t.Error("expected Ok flag true")
	}
	if Ok("x").Err().V {
		t.Error("expected Err flag false")
	}
	e := Err[string](errors.New("boom"))
	if e.Ok().V {
		t.Error("expected Ok flag false")
	}
	if !e.Err().V {
		t.Error("expected Err flag true")
	}
}

func TestResErrorAndGet(t *testing.T) {
	e := errors.New("oops")
	r := ResOf(5, e)
	if r.Error() != e {
		t.Errorf("Error mismatch: %v", r.Error())
	}
	v, err := r.Get()
	if v != 5 || err != e {
		t.Errorf("Get mismatch: %v %v", v, err)
	}
}

func TestResMust(t *testing.T) {
	ok := Ok(42)
	if ok.Must() != 42 {
		t.Errorf("Must wrong: %v", ok.Must())
	}
	err := errors.New("bad")
	mustPanicRes(t, func() { Err[int](err).Must() }, errorPrefixRes+err.Error())
}

func TestResMustErr(t *testing.T) {
	err := errors.New("abc")
	r := Err[int](err)
	if r.MustErr() != err {
		t.Errorf("MustErr wrong: %v", r.MustErr())
	}
	mustPanicRes(t, func() { _ = Ok(1).MustErr() }, errorPrefixRes+"no error")
}

func TestResOrAndOrZero(t *testing.T) {
	ok := Ok("ok")
	if got := ok.Or("alt"); got != "ok" {
		t.Errorf("Or wrong: %v", got)
	}
	if got := ok.OrZero(); got != "ok" {
		t.Errorf("OrZero wrong: %v", got)
	}
	bad := Err[string](errors.New("xx"))
	if got := bad.Or("alt"); got != "alt" {
		t.Errorf("Or wrong: %v", got)
	}
	if got := bad.OrZero(); got != "" {
		t.Errorf("OrZero wrong: %v", got)
	}
}
