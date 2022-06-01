package quickjs

// #include "bridge.h"
import "C"
import "unsafe"

type ClassID struct {
	value C.JSClassID
}

func NewClassID() ClassID {
	var classID ClassID
	C.JS_NewClassID(&classID.value)
	return classID
}

func (i ClassID) SetProto(c Context, proto Value) {
	C.JS_SetClassProto(c.ref, i.value, proto.dup().ref)
}

func (c Context) NewClassInstance(classID ClassID) Value {
	return c.valueFromRef(C.JS_NewObjectClass(c.ref, C.int(classID.value)))
}

type ClassFinalizer func(Runtime, Value)

type ClassDef struct {
	ClassName string
}

func (c ClassDef) c() (C.JSClassDef, func()) {
	className := C.CString(c.ClassName)
	free := func() {
		C.free(unsafe.Pointer(className))
	}
	return C.JSClassDef{class_name: className}, free
}

func NewClass(r Runtime, classID ClassID, classDef ClassDef) error {
	cClassDef, free := classDef.c()
	defer free()
	if errno := C.JS_NewClass(r.ref, classID.value, &cClassDef); errno < 0 {
		return Errno(errno)
	}
	return nil
}
