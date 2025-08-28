package shallow

import (
	"fmt"
	"reflect"
	"sync"
	"unsafe"
)

var known sync.Map

type eface struct {
	t uintptr // points to runtime._type
	d unsafe.Pointer
}

type key struct {
	fromType uintptr
	toType   uintptr
}

type entry struct {
	maps []mapping
}

type mapping struct {
	fromOffset uintptr
	toOffset   uintptr
	size       uintptr
}

func Map(from, to any) {
	eFrom := (*eface)(unsafe.Pointer(&from))
	eTo := (*eface)(unsafe.Pointer(&to))
	found, ok := known.Load(key{eFrom.t, eTo.t})
	var e entry
	if ok {
		e = found.(entry)
	} else {
		e = computeMapping(from, to)
		known.Store(key{eFrom.t, eTo.t}, e)
	}
	doMap(eFrom.d, eTo.d, e.maps)
}

func computeMapping(from, to any) entry {
	tFrom := reflect.TypeOf(from)
	tTo := reflect.TypeOf(to)
	if tFrom.Kind() == reflect.Ptr {
		tFrom = tFrom.Elem()
	} else if tFrom.Kind() != reflect.Struct {
		panic(fmt.Sprintf("from must be struct or *struct:\nfrom: %v\n  to: %v", tFrom, tTo))
	} else if tTo.Kind() != reflect.Pointer {
		panic(fmt.Sprintf("to must be pointer to struct:\nfrom: %v\n  to: %v", tFrom, tTo))
	}
	tTo = tTo.Elem()
	if tTo.Kind() != reflect.Struct {
		panic(fmt.Sprintf("to must be pointer to struct:\nfrom: %v\n  to: %v", tFrom, tTo))
	}
	var maps []mapping
	for i := 0; i < tTo.NumField(); i++ {
		df := tTo.Field(i)
		sf, ok := tFrom.FieldByName(df.Name)
		if !ok {
			panic(fmt.Sprintf("field '%s' not found in source:\nfrom:%v\n  to:  %v", df.Name, tFrom, tTo))
		} else if df.Type != sf.Type {
			panic(fmt.Sprintf("field '%s' has mismatched type:\nfrom:%T\n  to:%T", df.Name, sf.Type, df.Type))
		}
		maps = append(maps, mapping{
			fromOffset: sf.Offset,
			toOffset:   df.Offset,
			size:       df.Type.Size(),
		})
	}
	return entry{maps: maps}
}

func doMap(from, to unsafe.Pointer, mappings []mapping) {
	for _, m := range mappings {
		fromPtr := unsafe.Add(from, m.fromOffset)
		toPtr := unsafe.Add(to, m.toOffset)
		memmove(toPtr, fromPtr, m.size)
	}
}

func memmove(dst, src unsafe.Pointer, n uintptr) {
	var dstPtr = (*byte)(dst)
	var srcPtr = (*byte)(src)
	dstSlice := unsafe.Slice(dstPtr, n)
	srcSlice := unsafe.Slice(srcPtr, n)
	copy(dstSlice, srcSlice)
}
