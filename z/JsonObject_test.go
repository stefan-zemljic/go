package z

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"testing"
)

func JVNum(n string) JsonValue {
	return JsonValue{number(n)}
}

func JVBool(b bool) JsonValue {
	return JsonValue{b}
}

func JVStr(s string) JsonValue {
	return JsonValue{s}
}

func JVNull() JsonValue {
	return JsonValue{data: nil}
}

func newObj(pairs ...struct {
	K string
	V JsonValue
}) *JsonObject {
	o := &JsonObject{}
	for _, p := range pairs {
		o.Put(p.K, p.V)
	}
	return o
}

func TestLenKeyValue(t *testing.T) {
	o := newObj(
		struct {
			K string
			V JsonValue
		}{"a", JVNum("1")},
		struct {
			K string
			V JsonValue
		}{"b", JVNum("2")},
	)
	if o.Len() != 2 {
		t.Fatal("Len wrong")
	}
	if o.Key(0) != "a" || o.Value(1).data != number("2") {
		t.Fatal("Key/Value wrong")
	}
}

func TestHas(t *testing.T) {
	o := &JsonObject{
		keys:   []string{"a"},
		values: []JsonValue{JVNum("1")},
	}
	if o.Has("a") != True {
		t.Fatal("expected has a")
	}
	if o.Has("b") != False {
		t.Fatal("expected no b")
	}
}

func TestGet(t *testing.T) {
	o := &JsonObject{}
	if !o.Get("none").Empty().V {
		t.Fatal("expected empty get")
	}
	o.Put("x", JVNum("42"))
	if v := o.Get("x"); v.Empty().V || v.V.data != number("42") {
		t.Fatal("expected 42")
	}
}

func TestPutNewAndReplace(t *testing.T) {
	o := &JsonObject{}
	if !o.Put("x", JVStr("one")).Empty().V {
		t.Fatal("expected empty box when new key")
	}
	if old := o.Put("x", JVStr("two")); old.Empty().V || old.V.data != "one" {
		t.Fatal("expected old='one'")
	}
	if o.Value(0).data != "two" {
		t.Fatal("value not replaced")
	}
}

func TestDelete(t *testing.T) {
	o := newObj(
		struct {
			K string
			V JsonValue
		}{"a", JVBool(true)},
		struct {
			K string
			V JsonValue
		}{"b", JVBool(false)},
	)
	if o.Delete("z") != False {
		t.Fatal("expected false on missing")
	}
	if o.Delete("a") != True {
		t.Fatal("expected true on delete")
	}
	if len(o.keys) != 1 || o.keys[0] != "b" || o.indices["b"] != 0 {
		t.Fatal("delete did not reindex properly")
	}
}

func TestDeleteP(t *testing.T) {
	o := newObj(struct {
		K string
		V JsonValue
	}{"a", JVNum("1")})
	o.DeleteP("a")
	if o.Len() != 0 {
		t.Fatal("expected empty after deleteP")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		} else if !reflect.DeepEqual(r, fmt.Sprintf("%skey %q not found in object", errorPrefixObject, "nope")) {
			t.Fatalf("unexpected panic: %v", r)
		}
	}()
	o.DeleteP("nope")
}

func TestDeleteInitCase(t *testing.T) {
	o := &JsonObject{}
	if o.Delete("missing") != False {
		t.Fatal("expected false when deleting from empty uninitialized map")
	}
	if o.indices == nil {
		t.Fatal("expected mapd initialized")
	}
}

func TestClear(t *testing.T) {
	o := newObj(struct {
		K string
		V JsonValue
	}{"a", JVStr("x")})
	o.Clear()
	if o.Len() != 0 || len(o.indices) != 0 {
		t.Fatal("clear failed")
	}
}

func TestKeysValuesEntries(t *testing.T) {
	o := newObj(
		struct {
			K string
			V JsonValue
		}{"a", JVNum("1")},
		struct {
			K string
			V JsonValue
		}{"b", JVNum("2")},
	)
	ks := o.Keys()
	vs := o.Values()
	ens := o.Entries()
	if !reflect.DeepEqual(ks, []string{"a", "b"}) {
		t.Fatal("keys wrong")
	}
	if !reflect.DeepEqual(vs, []JsonValue{JVNum("1"), JVNum("2")}) {
		t.Fatal("vals wrong")
	}
	if len(ens) != 2 || ens[0].K != "a" || ens[1].V.data != number("2") {
		t.Fatal("entries wrong")
	}
	ks[0] = "zzz"
	vs[0] = JVStr("hack")
	if o.keys[0] == "zzz" || o.values[0].data == "hack" {
		t.Fatal("expected clone not alias")
	}
}

func TestJsonObject_MarshalJSON_ValueError(t *testing.T) {
	o := JsonObject{}
	o.keys = []string{"bad"}
	o.values = []JsonValue{{data: math.NaN()}}
	_, err := o.MarshalJSON()
	if err == nil || !strings.Contains(err.Error(), "value") {
		t.Fatalf("expected value marshal error, got: %v", err)
	}
}

func TestJsonObject_UnmarshalJSON_Errors(t *testing.T) {
	var o JsonObject
	if err := o.UnmarshalJSON([]byte(`{"a":}`)); err == nil {
		t.Errorf("expected syntax error")
	}
	if err := o.UnmarshalJSON([]byte(`123`)); err == nil {
		t.Errorf("expected 'object expected' error")
	}
	if err := o.UnmarshalJSON([]byte(`{"a": 1`)); err == nil {
		t.Errorf("expected truncated JSON error")
	}
	if err := o.UnmarshalJSON([]byte(`{"a":1]`)); err == nil {
		t.Errorf("expected wrong end token error")
	}
}

func TestJsonObject_UnmarshalJSON(t *testing.T) {
	json := `{"a":1,"b":[true,false,null],"c":{"d":"x"}}`
	var o JsonObject
	if err := o.UnmarshalJSON([]byte(json)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := o.MarshalJSON()
	if err != nil {
		t.Fatalf("unexpected marshal error: %v", err)
	} else if string(b) != json {
		t.Fatalf("marshal mismatch: got %s", b)
	}
}
