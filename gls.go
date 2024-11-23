package gls

import (
	"reflect"
	"sync"
	"unsafe"
)

// ID returns goroutine id of this goroutine
func ID() int64 {
	return idOf(getg(), goidoff)
}

// ParentID return goid of goroutine that created this goroutine
func ParentID() int64 {
	return idOf(getg(), parentgoidoff)
}

func idOf(g unsafe.Pointer, off uintptr) int64 {
	return *(*int64)(add(g, off))
}

//go:nosplit
func getg() unsafe.Pointer {
	return *(*unsafe.Pointer)(add(getm(), curgoff))
}

//go:linkname add runtime.add
//go:nosplit
func add(p unsafe.Pointer, x uintptr) unsafe.Pointer

//go:linkname getm runtime.getm
//go:nosplit
func getm() unsafe.Pointer

var (
	curgoff       = offset("*runtime.m", "curg")
	goidoff       = offset("*runtime.g", "goid")
	parentgoidoff = offset("*runtime.g", "parentGoid")
)

// offset returns the offset into typ for the given field.
func offset(typ, field string) uintptr {
	offsetOnce.Do(func() {
		offsetCache = make(map[string]uintptr)
	})

	offset, ok := offsetCache[typ+field]
	if !ok {
		rt := toType(typesByString(typ)[0])
		f, _ := rt.Elem().FieldByName(field)
		offset = f.Offset
		offsetCache[typ+field] = offset
	}

	return offset
}

//go:linkname typesByString reflect.typesByString
func typesByString(s string) []unsafe.Pointer

//go:linkname toType reflect.toType
func toType(t unsafe.Pointer) reflect.Type

var offsetOnce sync.Once
var offsetCache map[string]uintptr
