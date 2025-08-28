package env

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

func Read(target any) {
	pointer_ := reflect.ValueOf(target)
	if pointer_.Kind() != reflect.Pointer || pointer_.IsNil() {
		panic("target must be a non-nil pointer")
	}
	struct_ := pointer_.Elem()
	if struct_.Kind() != reflect.Struct {
		panic("target must point to a struct")
	}
	type_ := struct_.Type()
	for i := 0; i < type_.NumField(); i++ {
		field := struct_.Field(i)
		fieldType := type_.Field(i)
		name := fieldType.Name
		if !field.CanSet() {
			panic(fmt.Sprintf("field %q cannot be set", type_.Field(i).Name))
		}
		varName := upperSnakeCase(name)
		value, ok := os.LookupEnv(varName)
		if !ok {
			panic(fmt.Sprintf("environment variable %q not set", varName))
		}
		if field.Kind() != reflect.String {
			panic(fmt.Sprintf("field %q is not a string", name))
		}
	}
}

func upperSnakeCase(s string) string {
	result := ""
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result += "_"
		}
		result += string(r)
	}
	return strings.ToUpper(result)
}
