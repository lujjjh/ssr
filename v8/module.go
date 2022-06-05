package v8

// #include "v8_bridge.h"
import "C"
import "unsafe"

type Module struct {
	ptr C.v8_persistent_module_ptr
}

func (c *Isolate) CompileModule(source string, resourceName string) (*Module, error) {
	cSource := C.CString(source)
	defer C.free(unsafe.Pointer(cSource))

	cResourceName := C.CString(resourceName)
	defer C.free(unsafe.Pointer(cResourceName))

	// TODO: Handle errors.
	return &Module{C.v8_script_compiler_compile_module(c.ptr, cSource, cResourceName)}, nil
}

func (m *Module) Dispose() {
	C.v8_module_dispose(m.ptr)
}

func (m *Module) Run(c *Context) *Value {
	return valueFromPtr(c, C.v8_module_run(c.i.ptr, c.ptr, m.ptr))
}
