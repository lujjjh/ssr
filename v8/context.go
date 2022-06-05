package v8

// #include "v8_bridge.h"
import "C"

type Context struct {
	i   *Isolate
	ptr C.v8_persistent_context_ptr
}

func (i *Isolate) NewContext() *Context {
	return &Context{i, C.v8_context_new(i.ptr)}
}

func (c *Context) Dispose() {
	C.v8_context_dispose(c.ptr)
}

func (c *Context) Isolate() *Isolate {
	return c.i
}
