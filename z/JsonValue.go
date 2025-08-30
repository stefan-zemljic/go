package z

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

type JsonValue struct {
	data any
}

type JsonValueKind int

const (
	JsonValueKindNull JsonValueKind = iota
	JsonValueKindBool
	JsonValueKindNumber
	JsonValueKindString
	JsonValueKindArray
	JsonValueKindObject
)

type number string

func (s JsonValueKind) String() string {
	switch s {
	case JsonValueKindNull:
		return "null"
	case JsonValueKindBool:
		return "bool"
	case JsonValueKindNumber:
		return "number"
	case JsonValueKindString:
		return "string"
	case JsonValueKindArray:
		return "array"
	case JsonValueKindObject:
		return "object"
	default:
		panic(fmt.Sprintf("unknown JsonValueKind: %d", s))
	}
}

func (v JsonValue) Kind() JsonValueKind {
	switch v.data.(type) {
	case nil:
		return JsonValueKindNull
	case bool:
		return JsonValueKindBool
	case number:
		return JsonValueKindNumber
	case string:
		return JsonValueKindString
	case []JsonValue:
		return JsonValueKindArray
	case *JsonObject:
		return JsonValueKindObject
	default:
		panic(fmt.Sprintf("%sunknown value type: %T", errorPrefixJsonValue, v.data))
	}
}

func (v JsonValue) Null() Flag {
	return FlagOf(v.data == nil)
}

func (v JsonValue) Bool() Res[bool] {
	if b, ok := v.data.(bool); ok {
		return Ok(b)
	}
	return Errorf[bool]("%svalue is not a bool: %T", errorPrefixJsonValue, v.data)
}

func (v JsonValue) Int() Res[int] { return valueInt[int](v) }

func (v JsonValue) Int8() Res[int8] { return valueInt[int8](v) }

func (v JsonValue) Int16() Res[int16] { return valueInt[int16](v) }

func (v JsonValue) Int32() Res[int32] { return valueInt[int32](v) }

func (v JsonValue) Int64() Res[int64] { return valueInt[int64](v) }

func valueInt[T int | int8 | int16 | int32 | int64](v JsonValue) Res[T] {
	if n, ok := v.data.(number); ok {
		i, err := strconv.ParseInt(string(n), 10, 64)
		return ResOf[T](T(i), err)
	}
	return Errorf[T]("%svalue is not a number: %T", errorPrefixJsonValue, v.data)
}

func (v JsonValue) Uint() Res[uint] { return valueUint[uint](v) }

func (v JsonValue) Uint8() Res[uint8] { return valueUint[uint8](v) }

func (v JsonValue) Uint16() Res[uint16] { return valueUint[uint16](v) }

func (v JsonValue) Uint32() Res[uint32] { return valueUint[uint32](v) }

func (v JsonValue) Uint64() Res[uint64] { return valueUint[uint64](v) }

func (v JsonValue) Uintptr() Res[uintptr] { return valueUint[uintptr](v) }

func valueUint[T uint | uint8 | uint16 | uint32 | uint64 | uintptr](v JsonValue) Res[T] {
	if n, ok := v.data.(number); ok {
		u, err := strconv.ParseUint(string(n), 10, 64)
		if err != nil {
			return Errorf[T]("%sparsing number: %v", errorPrefixJsonValue, err)
		}
		return Ok(T(u))
	}
	return Errorf[T]("%svalue is not a number: %T", errorPrefixJsonValue, v.data)
}

func (v JsonValue) Float32() Res[float32] { return valueFloat[float32](v) }

func (v JsonValue) Float64() Res[float64] { return valueFloat[float64](v) }

func valueFloat[T float32 | float64](v JsonValue) Res[T] {
	if n, ok := v.data.(number); ok {
		f, err := strconv.ParseFloat(string(n), 64)
		return ResOf[T](T(f), err)
	}
	return Errorf[T]("%svalue is not a number: %T", errorPrefixJsonValue, v.data)
}

func (v JsonValue) Number() Res[string] {
	if n, ok := v.data.(number); ok {
		return Ok(string(n))
	}
	return Errorf[string]("%svalue is not a number: %T", errorPrefixJsonValue, v.data)
}

func (v JsonValue) String() Res[string] {
	if s, ok := v.data.(string); ok {
		return Ok(s)
	}
	return Errorf[string]("%svalue is not a string: %T", errorPrefixJsonValue, v.data)
}

func (v JsonValue) Array() Res[[]JsonValue] {
	if a, ok := v.data.([]JsonValue); ok {
		return Ok(a)
	}
	return Errorf[[]JsonValue]("%svalue is not an array: %T", errorPrefixJsonValue, v.data)
}

func (v JsonValue) Object() Res[*JsonObject] {
	if m, ok := v.data.(*JsonObject); ok {
		return Ok(m)
	}
	return Errorf[*JsonObject]("%svalue is not an object: %T", errorPrefixJsonValue, v.data)
}

func (v JsonValue) MarshalJSON() ([]byte, error) {
	switch d := v.data.(type) {
	case number:
		return Unwrap(string(d)), nil
	case []JsonValue:
		if d == nil {
			return []byte("[]"), nil
		}
		return json.Marshal(d)
	default:
		return json.Marshal(d)
	}
}

func (v *JsonValue) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	val, err := parseJson(dec)
	if err != nil {
		return err
	}
	*v = val
	return nil
}

const errorPrefixJsonValue = "z.JsonValue: "
