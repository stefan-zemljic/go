package z

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestJsonValue_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		json string
	}{
		{"null", `null`},
		{"bool_true", `true`},
		{"bool_false", `false`},
		{"number", `123`},
		{"string", `"hi"`},
		{"empty_array", `[]`},
		{"array_nested", `[1,"s",[],{"k":5}]`},
		{"empty_object", `{}`},
		{"object_simple", `{"a":1,"b":"x"}`},
		{"object_nested", `{"a":5,"b":[{"c":7},[]]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v JsonValue
			if err := json.Unmarshal([]byte(tt.json), &v); err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			out, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("Marshal error: %v", err)
			}
			if string(out) != tt.json {
				t.Fatalf("roundtrip mismatch:\n in : %s\nout: %s", tt.json, out)
			}
		})
	}
}

func TestJsonValueKindString(t *testing.T) {
	cases := []struct {
		kind JsonValueKind
		want string
	}{
		{JsonValueKindNull, "null"},
		{JsonValueKindBool, "bool"},
		{JsonValueKindNumber, "number"},
		{JsonValueKindString, "string"},
		{JsonValueKindArray, "array"},
		{JsonValueKindObject, "object"},
	}
	for _, c := range cases {
		if got := c.kind.String(); got != c.want {
			t.Errorf("got %q, want %q", got, c.want)
		}
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for unknown kind")
		}
	}()
	_ = JsonValueKind(999).String()
}

func TestJsonValueKind(t *testing.T) {
	cases := []struct {
		val  JsonValue
		kind JsonValueKind
	}{
		{JsonValue{nil}, JsonValueKindNull},
		{JsonValue{true}, JsonValueKindBool},
		{JsonValue{number("123")}, JsonValueKindNumber},
		{JsonValue{"str"}, JsonValueKindString},
		{JsonValue{[]JsonValue{}}, JsonValueKindArray},
		{JsonValue{&JsonObject{}}, JsonValueKindObject},
	}
	for _, c := range cases {
		if got := c.val.Kind(); got != c.kind {
			t.Errorf("expected %v, got %v", c.kind, got)
		}
	}
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for unknown type")
		}
	}()
	_ = JsonValue{42}.Kind()
}

func TestNullAndBool(t *testing.T) {
	if !(JsonValue{nil}.Null().V) {
		t.Errorf("expected true for null")
	}
	if (JsonValue{123}.Null().V) {
		t.Errorf("expected false for non-null")
	}
	ok, err := JsonValue{true}.Bool().Get()
	if err != nil || !ok {
		t.Errorf("expected true bool, got %v %v", ok, err)
	}
	_, err = JsonValue{"not-bool"}.Bool().Get()
	if err == nil {
		t.Errorf("expected error for non-bool")
	}
}

func TestIntParsing(t *testing.T) {
	v := JsonValue{number("42")}
	if got, err := v.Int().Get(); err != nil || got != 42 {
		t.Errorf("expected 42, got %v, err=%v", got, err)
	}
	v = JsonValue{"not-number"}
	_, err := v.Int().Get()
	if err == nil {
		t.Errorf("expected error for non-number")
	}
}

func TestUintParsing(t *testing.T) {
	v := JsonValue{number("42")}
	if got, err := v.Uint().Get(); err != nil || got != 42 {
		t.Errorf("expected 42, got %v, err=%v", got, err)
	}
	v = JsonValue{number("-1")}
	_, err := v.Uint().Get()
	if err == nil || !strings.Contains(err.Error(), "parsing number") {
		t.Errorf("expected parsing error, got %v", err)
	}
	v = JsonValue{"nope"}
	_, err = v.Uint().Get()
	if err == nil {
		t.Errorf("expected error for non-number")
	}
}

func TestFloatParsing(t *testing.T) {
	v := JsonValue{number("3.14")}
	if got, err := v.Float32().Get(); err != nil || got != float32(3.14) {
		t.Errorf("expected 3.14, got %v, err=%v", got, err)
	}
	v = JsonValue{"nope"}
	_, err := v.Float64().Get()
	if err == nil {
		t.Errorf("expected error for non-number")
	}
}

func TestNumberAndString(t *testing.T) {
	v := JsonValue{number("123")}
	if got, err := v.Number().Get(); err != nil || got != "123" {
		t.Errorf("expected 123, got %v, err=%v", got, err)
	}
	v = JsonValue{"str"}
	if got, err := v.String().Get(); err != nil || got != "str" {
		t.Errorf("expected str, got %v, err=%v", got, err)
	}
	_, err := JsonValue{123}.String().Get()
	if err == nil {
		t.Errorf("expected error for non-string")
	}
	_, err = JsonValue{"nope"}.Number().Get()
	if err == nil {
		t.Errorf("expected error for non-number in Number()")
	}
}

func TestArrayAndObject(t *testing.T) {
	arr := []JsonValue{{number("1")}}
	v := JsonValue{arr}
	got, err := v.Array().Get()
	if err != nil || !reflect.DeepEqual(got, arr) {
		t.Errorf("expected %v, got %v, err=%v", arr, got, err)
	}
	_, err = JsonValue{"not-array"}.Array().Get()
	if err == nil {
		t.Errorf("expected error for non-array")
	}
	obj := JsonObject{keys: []string{"a"}, values: []JsonValue{{number("1")}}}
	v = JsonValue{&obj}
	gotObj, err := v.Object().Get()
	if err != nil || !reflect.DeepEqual(gotObj, &obj) {
		t.Errorf("expected %v, got %v, err=%v", obj, gotObj, err)
	}
	_, err = JsonValue{"not-object"}.Object().Get()
	if err == nil {
		t.Errorf("expected error for non-object")
	}
}

func TestMarshalJSON(t *testing.T) {
	v := JsonValue{number("123")}
	b, _ := v.MarshalJSON()
	if string(b) != "123" {
		t.Errorf("expected 123, got %s", b)
	}
	v = JsonValue{[]JsonValue{}}
	b, _ = v.MarshalJSON()
	if string(b) != "[]" {
		t.Errorf("expected [], got %s", b)
	}
	v = JsonValue{[]JsonValue(nil)}
	b, _ = v.MarshalJSON()
	if string(b) != "[]" {
		t.Errorf("expected [], got %s", b)
	}
	obj := &JsonObject{keys: []string{"a"}, values: []JsonValue{{"b"}}}
	v = JsonValue{obj}
	b, _ = v.MarshalJSON()
	if string(b) != `{"a":"b"}` {
		t.Errorf("expected object json, got %s", b)
	}
	v = JsonValue{"str"}
	b, _ = v.MarshalJSON()
	if string(b) != `"str"` {
		t.Errorf("expected \"str\", got %s", b)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	var v JsonValue
	if err := v.UnmarshalJSON([]byte(`123`)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if k := v.Kind(); k != JsonValueKindNumber {
		t.Errorf("expected number kind, got %v", k)
	}
	var v2 JsonValue
	err := v2.UnmarshalJSON([]byte(`{invalid}`))
	if err == nil {
		t.Errorf("expected error for invalid json")
	}
}

func TestIntVariantParsingErrors(t *testing.T) {
	v := JsonValue{"oops"}
	if _, err := v.Int8().Get(); err == nil {
		t.Errorf("expected error for non-number in Int8")
	}
	if _, err := v.Int16().Get(); err == nil {
		t.Errorf("expected error for non-number in Int16")
	}
	if _, err := v.Int32().Get(); err == nil {
		t.Errorf("expected error for non-number in Int32")
	}
	if _, err := v.Int64().Get(); err == nil {
		t.Errorf("expected error for non-number in Int64")
	}
	v = JsonValue{number("12x")}
	if _, err := v.Int().Get(); err == nil {
		t.Errorf("expected parse error for invalid int")
	}
}

func TestUintVariantParsingErrors(t *testing.T) {
	v := JsonValue{number("-5")}
	if _, err := v.Uint8().Get(); err == nil {
		t.Errorf("expected parse error for negative number in Uint8")
	}
	if _, err := v.Uint16().Get(); err == nil {
		t.Errorf("expected parse error for negative number in Uint16")
	}
	v = JsonValue{"nope"}
	if _, err := v.Uint32().Get(); err == nil {
		t.Errorf("expected error for non-number in Uint32")
	}
	if _, err := v.Uint64().Get(); err == nil {
		t.Errorf("expected error for non-number in Uint64")
	}
	if _, err := v.Uintptr().Get(); err == nil {
		t.Errorf("expected error for non-number in Uintptr")
	}
}

func TestFloatParsingErrors(t *testing.T) {
	v := JsonValue{number("3.1.4")}
	if _, err := v.Float32().Get(); err == nil {
		t.Errorf("expected parse error for invalid float")
	}
}
