package v8

// #include "v8_bridge.h"
import "C"
import "unsafe"

type Value struct {
	c   *Context
	ptr C.v8_persistent_value_ptr
}

func valueFromPtr(c *Context, ptr C.v8_persistent_value_ptr) *Value {
	return &Value{c, ptr}
}

func (v *Value) Dispose() {
	C.v8_value_dispose(v.ptr)
}

func (v *Value) String() string {
	cStr := C.v8_value_to_detail_string(v.c.i.ptr, v.c.ptr, v.ptr)
	defer C.free(unsafe.Pointer(cStr))

	return C.GoString(cStr)
}

func (v *Value) Call() *Value {
	return valueFromPtr(v.c, C.v8_function_call(v.c.i.ptr, v.c.ptr, v.ptr))
}
