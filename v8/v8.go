package v8

// #cgo CXXFLAGS: -fno-rtti -fPIC -std=c++14 -DV8_COMPRESS_POINTERS -DV8_31BIT_SMIS_ON_64BIT_ARCH -I${SRCDIR}/include -Wall
// #cgo LDFLAGS: -pthread -lv8
// #cgo darwin LDFLAGS: -L${SRCDIR}/bin/darwin
// #cgo linux LDFLAGS: -L${SRCDIR}/bin/linux -ldl
//
// #include "v8_bridge.h"
import "C"

func Initialize() {
	C.v8_initialize()
}

func Dispose() {
	C.v8_dispose()
}

func GetVersion() string {
	return C.GoString(C.v8_get_version())
}
