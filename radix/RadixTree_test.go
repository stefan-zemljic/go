package radix

import (
	"reflect"
	"testing"
)

func newTree[V any]() *Tree[V] {
	return New[V]()
}
func TestLen(t *testing.T) {
	tr := newTree[int]()
	if tr.Len() != 0 {
		t.Errorf("expected 0, got %d", tr.Len())
	}
	tr.Insert("a", 1)
	if tr.Len() != 1 {
		t.Errorf("expected 1, got %d", tr.Len())
	}
}
func TestInsert(t *testing.T) {
	tr := newTree[int]()
	old, ok := tr.Insert("a", 1)
	if ok || old != 0 {
		t.Errorf("expected zero,false got %v,%v", old, ok)
	}
	old, ok = tr.Insert("a", 2)
	if !ok || old != 1 {
		t.Errorf("expected 1,true got %v,%v", old, ok)
	}
	val, _ := tr.Get("a")
	if val != 2 {
		t.Errorf("expected 2, got %v", val)
	}
}
func TestDelete(t *testing.T) {
	tr := newTree[int]()
	tr.Insert("a", 1)
	val, ok := tr.Delete("a")
	if !ok || val != 1 {
		t.Errorf("expected 1,true got %v,%v", val, ok)
	}
	val, ok = tr.Delete("missing")
	if ok || val != 0 {
		t.Errorf("expected zero,false got %v,%v", val, ok)
	}
}
func TestDeletePrefix(t *testing.T) {
	tr := newTree[int]()
	tr.Insert("abc", 1)
	tr.Insert("abd", 2)
	tr.Insert("xyz", 3)
	n := tr.DeletePrefix("ab")
	if n != 2 {
		t.Errorf("expected 2 deletions, got %d", n)
	}
	if tr.Len() != 1 {
		t.Errorf("expected 1 left, got %d", tr.Len())
	}
}
func TestGet(t *testing.T) {
	tr := newTree[int]()
	tr.Insert("a", 1)
	val, ok := tr.Get("a")
	if !ok || val != 1 {
		t.Errorf("expected 1,true got %v,%v", val, ok)
	}
	val, ok = tr.Get("missing")
	if ok || val != 0 {
		t.Errorf("expected zero,false got %v,%v", val, ok)
	}
}
func TestLongestPrefix(t *testing.T) {
	tr := newTree[int]()
	tr.Insert("foo", 1)
	tr.Insert("foobar", 2)
	k, v, ok := tr.LongestPrefix("foobaz")
	if !ok || k != "foo" || v != 1 {
		t.Errorf("expected foo->1,true got %v,%v,%v", k, v, ok)
	}
	k, v, ok = tr.LongestPrefix("zzz")
	if ok || k != "" || v != 0 {
		t.Errorf("expected empty,false got %v,%v,%v", k, v, ok)
	}
}
func TestMinimumMaximum(t *testing.T) {
	tr := newTree[int]()
	k, v, ok := tr.Minimum()
	if ok || k != "" || v != 0 {
		t.Errorf("expected empty,false got %v,%v,%v", k, v, ok)
	}
	k, v, ok = tr.Maximum()
	if ok || k != "" || v != 0 {
		t.Errorf("expected empty,false got %v,%v,%v", k, v, ok)
	}
	tr.Insert("b", 2)
	tr.Insert("a", 1)
	tr.Insert("c", 3)
	k, v, ok = tr.Minimum()
	if !ok || k != "a" || v != 1 {
		t.Errorf("expected a->1 got %v,%v,%v", k, v, ok)
	}
	k, v, ok = tr.Maximum()
	if !ok || k != "c" || v != 3 {
		t.Errorf("expected c->3 got %v,%v,%v", k, v, ok)
	}
}
func TestWalk(t *testing.T) {
	tr := newTree[int]()
	tr.Insert("a", 1)
	tr.Insert("b", 2)
	var collected []string
	for s := range tr.Walk {
		collected = append(collected, s)
	}
	if len(collected) != 2 {
		t.Errorf("expected 2 walked, got %d", len(collected))
	}
}
func TestWalkPrefix(t *testing.T) {
	tr := newTree[int]()
	tr.Insert("foo", 1)
	tr.Insert("foobar", 2)
	tr.Insert("bar", 3)
	var collected []string
	tr.WalkPrefix("foo", func(s string, v int) bool {
		collected = append(collected, s)
		return true
	})
	expected := map[string]bool{"foo": true, "foobar": true}
	if len(collected) != 2 {
		t.Errorf("expected 2 items, got %v", collected)
	}
	for _, k := range collected {
		if !expected[k] {
			t.Errorf("unexpected key in prefix walk: %s", k)
		}
	}
}
func TestWalkPath(t *testing.T) {
	tr := newTree[int]()
	tr.Insert("foo", 1)
	tr.Insert("foobar", 2)
	var collected []string
	tr.WalkPath("foobar", func(s string, v int) bool {
		collected = append(collected, s)
		return true
	})
	if !reflect.DeepEqual(collected, []string{"foo", "foobar"}) {
		t.Errorf("unexpected path walk: %v", collected)
	}
}
func TestToMap(t *testing.T) {
	tr := newTree[int]()
	tr.Insert("x", 100)
	tr.Insert("y", 200)
	m := tr.ToMap()
	if len(m) != 2 || m["x"] != 100 || m["y"] != 200 {
		t.Errorf("unexpected map: %#v", m)
	}
}
