package jso

import (
	"reflect"
	"testing"
)

func TestMultiMap_PutGetDeleteClear(t *testing.T) {
	var m MultiMap[string, int]
	m.Put("a", 1)
	if got := m.Get("a"); !reflect.DeepEqual(got, []int{1}) {
		t.Fatalf("expected [1], got %v", got)
	}
	m.Put("a", 2)
	if got := m.Get("a"); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Fatalf("expected [1,2], got %v", got)
	}
	m.Put("b", 3)
	if got := m.Get("b"); !reflect.DeepEqual(got, []int{3}) {
		t.Fatalf("expected [3], got %v", got)
	}
	m.Delete("a")
	if m.Get("a") != nil {
		t.Fatalf("expected nil after delete")
	}
	m.Delete("zzz")
	m.Clear()
	if !m.IsEmpty() || m.KeyCount() != 0 || m.ValueCount() != 0 {
		t.Fatalf("expected empty after clear")
	}
	m.Clear()
}
func TestMultiMap_KeysValuesEntriesMapClone(t *testing.T) {
	var m MultiMap[string, int]
	m.Put("x", 1)
	m.Put("x", 2)
	m.Put("y", 3)
	keys := m.Keys()
	values := m.Values()
	entries := m.Entries()
	cmap := m.Map()
	clone := m.Clone()
	if len(keys) == 0 || len(values) == 0 || len(entries) == 0 || len(cmap) == 0 {
		t.Fatalf("expected all collections non-empty")
	}
	cmap["x"][0] = 999
	if m.Get("x")[0] == 999 {
		t.Fatalf("Map did not deep clone slice")
	}
	if !reflect.DeepEqual(m.entries, clone.entries) {
		t.Fatalf("entries mismatch after clone")
	}
	if !reflect.DeepEqual(m.values, clone.values) {
		t.Fatalf("values mismatch after clone")
	}
	m.Put("z", 9)
	if _, ok := clone.values["z"]; ok {
		t.Fatalf("clone mutated with original")
	}
}
func TestMultiMap_IterJsonPrettyString(t *testing.T) {
	var m MultiMap[string, int]
	m.Put("a", 1)
	m.Put("b", 2)
	m.Put("a", 3)
	count := 0
	m.Iter(func(k string, v int) bool {
		count++
		return true
	})
	if count != m.ValueCount() {
		t.Fatalf("expected %d iterations, got %d", m.ValueCount(), count)
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
	m.JsonTo(&w)
	if w.String() == "" {
		t.Fatalf("expected non-empty Writer output")
	}
	if m.Json() == "" || m.Pretty() == "" || m.String() == "" {
		t.Fatalf("expected Json/Pretty/String non-empty")
	}
}
