package quickjs

// #include "bridge.h"
import "C"

type Runtime struct {
	ref *C.JSRuntime
}

func NewRuntime() Runtime {
	return Runtime{C.JS_NewRuntime()}
}

func (r Runtime) Free() {
	C.JS_FreeRuntime(r.ref)
}
