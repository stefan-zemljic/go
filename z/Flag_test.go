package z

import "testing"

func mustPanicFlag(t *testing.T, fn func(), want string) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			got := r.(string)
			if got != want {
				t.Fatalf("panic mismatch: got %q, want %q", got, want)
			}
		} else {
			t.Fatalf("expected panic %q, got none", want)
		}
	}()
	fn()
}

func TestFlagOfAndV(t *testing.T) {
	f1 := FlagOf(true)
	if !f1.V {
		t.Error("FlagOf(true) wrong value")
	}
	f2 := FlagOf(false)
	if f2.V {
		t.Error("FlagOf(false) wrong value")
	}
}

func TestFlagNeg(t *testing.T) {
	f := FlagOf(true)
	n := f.Neg()
	if n.V != false {
		t.Errorf("Neg mismatch: got %v", n.V)
	}
	if f.Neg().Neg().V != f.V {
		t.Error("double Neg failed")
	}
}

func TestFlagNot(t *testing.T) {
	if got := FlagOf(true).Not(); got {
		t.Errorf("Not mismatch: got %v", got)
	}
	if got := FlagOf(false).Not(); !got {
		t.Errorf("Not mismatch: got %v", got)
	}
}

func TestFlagMust(t *testing.T) {
	if got := FlagOf(true).Must(); !got {
		t.Errorf("Must mismatch: got %v", got)
	}
	mustPanicFlag(t, func() { FlagOf(false).Must() }, errorPrefixFlag+"flag not set")
}

func TestFlagMustNot(t *testing.T) {
	if got := FlagOf(false).MustNot(); got != true {
		t.Errorf("MustNot mismatch: got %v", got)
	}
	mustPanicFlag(t, func() { FlagOf(true).MustNot() }, errorPrefixFlag+"flag set")
}
