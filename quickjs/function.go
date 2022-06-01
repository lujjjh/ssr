package quickjs

// #include "bridge.h"
import "C"
import (
	"unsafe"
)

type GoFunction func(c Context, this Value, args ...Value) Value

var goFunctions = make(map[*C.JSContext][]GoFunction)

func (c Context) freeGoFunctions() {
	delete(goFunctions, c.ref)
}

//export cFunctionProxy
func cFunctionProxy(ctx *C.JSContext, this C.JSValueConst, argc C.int, argv *C.JSValueConst, magic C.int, data *C.JSValue) C.JSValue {
	c := contextFromRef(ctx)
	refs := unsafe.Slice(argv, int(argc))
	args := make([]Value, argc)
	for i := 0; i < int(argc); i++ {
		args[i] = c.valueFromRef(refs[i])
	}
	dataRefs := unsafe.Slice(data, 1)
	id := c.valueFromRef(dataRefs[0]).GetInt32()
	f := goFunctions[c.ref][int(id)]
	return f(contextFromRef(ctx), c.valueFromRef(this), args...).ref
}

func (c Context) NewGoFunction(f GoFunction, length int, magic int) Value {
	goFunctions[c.ref] = append(goFunctions[c.ref], f)
	id := C.JS_NewInt32(c.ref, C.int(len(goFunctions[c.ref])-1))
	return c.valueFromRef(C.JS_NewCFunctionData(c.ref, (*C.JSCFunctionData)(unsafe.Pointer(C.cFunctionProxy)),
		C.int(length), C.int(magic), C.int(1), &id))
}

func (v Value) Call(this Value, args ...Value) Value {
	c := v.c
	if len(args) == 0 {
		return c.valueFromRef(C.JS_Call(c.ref, v.ref, this.ref, C.int(0), nil))
	}
	refs := make([]C.JSValue, len(args))
	for i := range args {
		refs[i] = args[i].ref
	}
	return c.valueFromRef(C.JS_Call(c.ref, v.ref, this.ref, C.int(len(args)), &refs[0]))
}
