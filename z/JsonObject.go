package z

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"
)

type JsonObject struct {
	keys    []string
	values  []JsonValue
	indices map[string]int
}

func (s *JsonObject) Len() int {
	return len(s.keys)
}

func (s *JsonObject) Key(i int) string {
	return s.keys[i]
}

func (s *JsonObject) Value(i int) JsonValue {
	return s.values[i]
}

func (s *JsonObject) Has(key string) Flag {
	if s.indices == nil {
		s.initMap()
	}
	_, ok := s.indices[key]
	return FlagOf(ok)
}

func (s *JsonObject) Get(key string) Box[JsonValue] {
	if s.indices == nil {
		s.initMap()
	}
	if i, ok := s.indices[key]; ok {
		return Just(s.values[i])
	}
	return Empty[JsonValue]()
}

func (s *JsonObject) Put(key string, val JsonValue) Box[JsonValue] {
	if s.indices == nil {
		s.initMap()
	}
	if i, ok := s.indices[key]; ok {
		old := s.values[i]
		s.values[i] = val
		return Just(old)
	}
	s.keys = append(s.keys, key)
	s.values = append(s.values, val)
	s.indices[key] = len(s.keys) - 1
	return Empty[JsonValue]()
}

func (s *JsonObject) Delete(key string) Flag {
	if s.indices == nil {
		s.initMap()
	}
	i, ok := s.indices[key]
	if !ok {
		return False
	}
	delete(s.indices, key)
	s.keys = append(s.keys[:i], s.keys[i+1:]...)
	s.values = append(s.values[:i], s.values[i+1:]...)
	for j := i; j < len(s.keys); j++ {
		s.indices[s.keys[j]] = j
	}
	return True
}

func (s *JsonObject) DeleteP(key string) {
	if !s.Delete(key).V {
		panic(fmt.Sprintf("%skey %q not found in object", errorPrefixObject, key))
	}
}

func (s *JsonObject) Clear() {
	s.keys = s.keys[:0]
	s.values = s.values[:0]
	s.indices = map[string]int{}
}

func (s *JsonObject) Keys() []string {
	return slices.Clone(s.keys)
}

func (s *JsonObject) Values() []JsonValue {
	return slices.Clone(s.values)
}

func (s *JsonObject) Entries() []struct {
	K string
	V JsonValue
} {
	res := make([]struct {
		K string
		V JsonValue
	}, len(s.keys))
	for i, k := range s.keys {
		res[i].K = k
		res[i].V = s.values[i]
	}
	return res
}

func (s *JsonObject) MarshalJSON() ([]byte, error) {
	var b strings.Builder
	b.WriteByte('{')
	for i, k := range s.keys {
		if i > 0 {
			b.WriteByte(',')
		}
		bk, _ := json.Marshal(k)
		b.Write(bk)
		b.WriteByte(':')
		bv, err := json.Marshal(s.values[i])
		if err != nil {
			return nil, fmt.Errorf("%svalue for key %q (%d) marshal error: %w", errorPrefixObject, k, i, err)
		}
		b.Write(bv)
	}
	b.WriteByte('}')
	return Unwrap(b.String()), nil
}

func (s *JsonObject) UnmarshalJSON(data []byte) error {
	v, err := ParseJson(data)
	if err != nil {
		return err
	}
	obj, ok := v.data.(*JsonObject)
	if !ok {
		return fmt.Errorf("%sobject expected, got %T", errorPrefixObject, v.data)
	}
	*s = *obj
	return nil
}

func (s *JsonObject) initMap() {
	s.indices = make(map[string]int, len(s.keys))
	for i, k := range s.keys {
		s.indices[k] = i
	}
}

const errorPrefixObject = "z.Object: "
