package jso

import (
	"reflect"
	"testing"
)

func TestMap_PutGetMustGetDelete(t *testing.T) {
	var m Map[string, int]
	m.Put("a", 1)
	if v, ok := m.Get("a"); !ok || v != 1 {
		t.Fatalf("expected 1, got %v %v", v, ok)
	}
	if v := m.MustGet("a"); v != 1 {
		t.Fatalf("expected MustGet=1, got %v", v)
	}
	m.Put("a", 2)
	if v := m.MustGet("a"); v != 2 {
		t.Fatalf("expected 2 after overwrite, got %v", v)
	}
	m.Put("b", 3)
	if v, ok := m.Get("b"); !ok || v != 3 {
		t.Fatalf("expected 3, got %v %v", v, ok)
	}
	m.Delete("a")
	if _, ok := m.Get("a"); ok {
		t.Fatalf("expected 'a' deleted")
	}
	m.Delete("zzz")
	m.Clear()
	if !m.IsEmpty() || m.Len() != 0 {
		t.Fatalf("expected empty after Clear")
	}
	if v := m.MustGet("nope"); v != 0 {
		t.Fatalf("expected 0, got %v", v)
	}
}
func TestMap_KeysValuesEntriesMapClone(t *testing.T) {
	var m Map[string, int]
	m.Put("x", 10)
	m.Put("y", 20)
	keys := m.Keys()
	values := m.Values()
	entries := m.Entries()
	cmap := m.Map()
	clone := m.Clone()
	if len(keys) != 2 || len(values) != 2 || len(entries) != 2 || len(cmap) != 2 {
		t.Fatalf("expected all collections size=2")
	}
	if !reflect.DeepEqual(m.entries, clone.entries) {
		t.Fatalf("entries mismatch after clone")
	}
	if !reflect.DeepEqual(m.values, clone.values) {
		t.Fatalf("values mismatch after clone")
	}
	m.Put("z", 30)
	if _, ok := clone.values["z"]; ok {
		t.Fatalf("clone mutated with original")
	}
}
func TestMap_IterToWriterToPrettyString(t *testing.T) {
	var m Map[string, int]
	m.Put("a", 1)
	m.Put("b", 2)
	count := 0
	m.Iter(func(k string, v int) bool {
		count++
		return true
	})
	if count != m.Len() {
		t.Fatalf("expected %d iterations, got %d", m.Len(), count)
	}
	count = 0
	m.Iter(func(k string, v int) bool {
		count++
		return false
	})
	if count != 1 {
		t.Fatalf("expected break after 1 iteration, got %d", count)
	}
	var w Writer
	m.ToWriter(&w)
	if w.String() == "" {
		t.Fatalf("expected non-empty String() from Writer")
	}
	if m.Pretty() == "" || m.String() == "" {
		t.Fatalf("expected Pretty and String not empty")
	}
	var st Stream
	m.To(&st)
	if len(st.Tokens) == 0 {
		t.Fatalf("expected Stream items")
	}
}
