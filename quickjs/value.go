package quickjs

// #include "bridge.h"
import "C"
import "unsafe"

type Value struct {
	c   Context
	ref C.JSValue
}

func (v Value) Free() {
	C.JS_FreeValue(v.c.ref, v.ref)
}

func (v Value) String() string {
	cString := C.JS_ToCString(v.c.ref, v.ref)
	defer C.JS_FreeCString(v.c.ref, cString)
	return C.GoString(cString)
}

func (v Value) IsException() bool {
	return C.JS_IsException(v.ref) != 0
}

func (v Value) GetInt32() int32 {
	return int32(C.JSValue_GetInt32(v.ref))
}

func (v Value) dup() Value {
	return v.c.valueFromRef(C.JS_DupValue(v.c.ref, v.ref))
}

type PropFlags C.int

const (
	PropConfigurable PropFlags = C.JS_PROP_CONFIGURABLE
	PropWritable     PropFlags = C.JS_PROP_WRITABLE
	PropEnumerable   PropFlags = C.JS_PROP_ENUMERABLE

	PropNormal PropFlags = C.JS_PROP_NORMAL
	PropGetSet PropFlags = C.JS_PROP_GETSET
)

func (v Value) DefinePropertyValue(prop string, value Value, flags PropFlags) error {
	cProp := C.CString(prop)
	defer C.free(unsafe.Pointer(cProp))
	errno := C.JS_DefinePropertyValueStr(v.c.ref, v.ref, cProp, value.dup().ref, C.int(flags))
	if errno != 0 {
		return Errno(errno)
	}
	return nil
}

func (v Value) SetConstructorBit(value bool) {
	if value {
		C.JS_SetConstructorBit(v.c.ref, v.ref, C.int(1))
	} else {
		C.JS_SetConstructorBit(v.c.ref, v.ref, C.int(0))
	}
}

func (v Value) SetConstructor(value Value) {
	C.JS_SetConstructor(v.c.ref, v.ref, value.ref)
}

func (c Context) valueFromRef(ref C.JSValue) Value {
	return Value{c: c, ref: ref}
}

func (c Context) Null() Value {
	return c.valueFromRef(C.JS_NULL)
}

func (c Context) Undefined() Value {
	return c.valueFromRef(C.JS_UNDEFINED)
}

func (c Context) False() Value {
	return c.valueFromRef(C.JS_FALSE)
}

func (c Context) True() Value {
	return c.valueFromRef(C.JS_TRUE)
}

func (c Context) Exception() Value {
	return c.valueFromRef(C.JS_EXCEPTION)
}

func (c Context) Uninitialized() Value {
	return c.valueFromRef(C.JS_UNINITIALIZED)
}
