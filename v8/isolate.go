package v8

// #include "v8_bridge.h"
import "C"

type Isolate struct {
	ptr C.v8_isolate_ptr
}

func NewIsolate() *Isolate {
	return &Isolate{C.v8_isolate_new()}
}

func (i *Isolate) Dispose() {
	C.v8_isolate_dispose(i.ptr)
}
