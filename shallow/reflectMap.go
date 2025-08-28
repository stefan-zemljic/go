package shallow

import (
	"fmt"
	"reflect"
)

func reflectMap(from, to any) {
	tFrom := reflect.TypeOf(from)
	tTo := reflect.TypeOf(to)
	if tFrom.Kind() == reflect.Ptr {
		tFrom = tFrom.Elem()
	} else if tFrom.Kind() != reflect.Struct {
		panic(fmt.Errorf("from must be struct or *struct:\nfrom: %v\n  to: %v", tFrom, tTo))
	} else if tTo.Kind() != reflect.Pointer {
		panic(fmt.Errorf("to must be pointer to struct:\nfrom: %v\n  to: %v", tFrom, tTo))
	}
	tToElem := tTo.Elem()
	if tToElem.Kind() != reflect.Struct {
		panic(fmt.Errorf("to must be pointer to struct:\nfrom: %v\n  to: %v", tFrom, tTo))
	}
	vFrom := reflect.ValueOf(from)
	if vFrom.Kind() == reflect.Pointer {
		vFrom = vFrom.Elem()
	}
	vTo := reflect.ValueOf(to).Elem()
	for i := 0; i < tToElem.NumField(); i++ {
		df := tToElem.Field(i)
		tf := vTo.Field(i)

		sf, ok := tFrom.FieldByName(df.Name)
		if !ok {
			panic(fmt.Errorf("field '%s' not found in source:\nfrom:%v\n  to:  %v", df.Name, tFrom, tToElem))
		} else if df.Type != sf.Type {
			panic(fmt.Errorf("field '%s' has mismatched type:\nfrom:%T\n  to:%T", df.Name, sf.Type, df.Type))
		}
		tf.Set(vFrom.FieldByIndex(sf.Index))
	}
}
