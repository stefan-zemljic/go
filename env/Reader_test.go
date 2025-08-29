package env

import (
	"os"
	"reflect"
	"testing"
)

type good struct {
	Name string
}
type notString struct {
	Name int
}
type unexported struct {
	name string
}

func TestRead_Success(t *testing.T) {
	err := os.Setenv("NAME", "hello")
	if err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}
	defer func() {
		err = os.Unsetenv("NAME")
		if err != nil {
			t.Fatalf("failed to unset env var: %v", err)
		}
	}()
	var g good
	Read(&g)
	if g.Name != "hello" {
		t.Errorf("expected hello, got %q", g.Name)
	}
}
func TestRead_TargetNotPointerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	var g good
	Read(g) // not a pointer
}
func TestRead_TargetNilPointerPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	var g *good
	Read(g)
}

func TestRead_TargetNotStructPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	x := 42
	Read(&x)
}
func TestRead_FieldCannotBeSetPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	var u unexported
	Read(&u)
}
func TestRead_EnvVarNotSetPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	var g good
	Read(&g)
}
func TestRead_FieldNotStringPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic, got none")
		}
	}()
	err := os.Setenv("NAME", "hello")
	if err != nil {
		t.Fatalf("failed to set env var: %v", err)
	}
	var ns notString
	Read(&ns)
}
func TestUpperSnakeCase(t *testing.T) {
	cases := map[string]string{
		"Name":    "NAME",
		"UserId":  "USER_ID",
		"ABC":     "A_B_C",
		"already": "ALREADY",
	}
	for in, want := range cases {
		got := upperSnakeCase(in)
		if got != want {
			t.Errorf("upperSnakeCase(%q) = %q, want %q", in, got, want)
		}
	}
}
func TestUpperSnakeCaseUnderscoreInsertion(t *testing.T) {
	got := upperSnakeCase("GoLang")
	want := "GO_LANG"
	if got != want {
		t.Errorf("expected %q, got %q", want, got)
	}
}
func TestUpperSnakeCaseEmpty(t *testing.T) {
	if got := upperSnakeCase(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}
func TestReadTypeSignature(t *testing.T) {
	typ := reflect.TypeOf(Read)
	if typ.Kind() != reflect.Func {
		t.Errorf("expected function, got %v", typ.Kind())
	}
}
