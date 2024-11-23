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
		offsetCache = cache{
			offsets: make(map[string]uintptr),
		}
	})

	offset, ok := offsetCache.get(typ + field)
	if !ok {
		rt := toType(typesByString(typ)[0])
		f, _ := rt.Elem().FieldByName(field)
		offset = f.Offset
		offsetCache.set(typ+field, offset)
	}

	return offset
}

//go:linkname typesByString reflect.typesByString
func typesByString(s string) []unsafe.Pointer

//go:linkname toType reflect.toType
func toType(t unsafe.Pointer) reflect.Type

var offsetOnce sync.Once
var offsetCache cache

type cache struct {
	mu      sync.RWMutex
	offsets map[string]uintptr
}

func (c *cache) get(key string) (uintptr, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.offsets[key]
	return v, ok
}

func (c *cache) set(key string, v uintptr) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.offsets[key] = v
}
