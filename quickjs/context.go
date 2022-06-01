package quickjs

// #include "bridge.h"
import "C"
import "unsafe"

type Context struct {
	ref *C.JSContext
}

func contextFromRef(ref *C.JSContext) Context {
	return Context{ref}
}

func (r Runtime) NewContext() Context {
	return contextFromRef(C.JS_NewContext(r.ref))
}

func (c Context) Free() {
	c.freeGoFunctions()
	C.JS_FreeContext(c.ref)
}

type EvalFlags int

var (
	EvalTypeGlobal           = EvalFlags(C.JS_EVAL_TYPE_GLOBAL)
	EvalTypeModule           = EvalFlags(C.JS_EVAL_TYPE_MODULE)
	EvalFlagStrict           = EvalFlags(C.JS_EVAL_FLAG_STRICT)
	EvalFlagStrip            = EvalFlags(C.JS_EVAL_FLAG_STRIP)
	EvalFlagCompileOnly      = EvalFlags(C.JS_EVAL_FLAG_COMPILE_ONLY)
	EvalFlagBacktraceBarrier = EvalFlags(C.JS_EVAL_FLAG_BACKTRACE_BARRIER)
)

func (c Context) Eval(input string, filename string, flags EvalFlags) Value {
	cInput := C.CString(input)
	defer C.free(unsafe.Pointer(cInput))

	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cInputLen := C.strlen(cInput)

	return c.valueFromRef(C.JS_Eval(c.ref, cInput, cInputLen, cFilename, C.int(flags)))
}

func (c Context) EvalFunction(f Value) Value {
	return c.valueFromRef(C.JS_EvalFunction(c.ref, f.dup().ref))
}

func (c Context) GetModuleExport(m Value, exportName string) Value {
	cExportName := C.CString(exportName)
	defer C.free(unsafe.Pointer(cExportName))

	return c.valueFromRef(C.JS_GetModuleExport(c.ref, C.JS_GetModuleDef(m.ref), cExportName))
}

func (c Context) GetException() Value {
	return c.valueFromRef(C.JS_GetException(c.ref))
}

func (c Context) GetGlobalObject() Value {
	value := c.valueFromRef(C.JS_GetGlobalObject(c.ref))
	value.Free()
	return value
}
