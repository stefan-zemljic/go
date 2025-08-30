package z

import (
	"testing"
)

func mustPanic(t *testing.T, fn func(), want string) {
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
