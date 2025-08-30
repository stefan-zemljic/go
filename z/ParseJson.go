package z

import (
	"bytes"
	"encoding/json"
)

func ParseJson(bs []byte) (JsonValue, error) {
	dec := json.NewDecoder(bytes.NewReader(bs))
	dec.UseNumber()
	return parseJson(dec)
}

func parseJson(dec *json.Decoder) (JsonValue, error) {
	tok, err := dec.Token()
	if err != nil {
		return JsonValue{}, err
	}
	switch t := tok.(type) {
	default:
		return JsonValue{nil}, nil
	case bool:
		return JsonValue{t}, nil
	case string:
		return JsonValue{t}, nil
	case json.Number:
		return JsonValue{number(t.String())}, nil
	case json.Delim:
		switch t {
		case '{':
			var ktok json.Token
			var val JsonValue
			obj := JsonObject{}
			obj.Clear()
			for dec.More() {
				ktok, _ = dec.Token()
				val, err = parseJson(dec)
				if err != nil {
					return JsonValue{}, err
				}
				obj.keys = append(obj.keys, ktok.(string))
				obj.values = append(obj.values, val)
			}
			if _, err = dec.Token(); err != nil {
				return JsonValue{}, err
			}
			return JsonValue{&obj}, nil
		default:
			var arr []JsonValue
			var elem JsonValue
			for dec.More() {
				elem, err = parseJson(dec)
				if err != nil {
					return JsonValue{}, err
				}
				arr = append(arr, elem)
			}
			if _, err = dec.Token(); err != nil {
				return JsonValue{}, err
			}
			if arr == nil {
				arr = make([]JsonValue, 0)
			}
			return JsonValue{arr}, nil
		}
	}
}
