package g

/*
#cgo pkg-config: glib-2.0
#include <glib.h>
*/
import "C"

import (
	"structs"
	"unsafe"
)

type Error struct {
	_ structs.HostLayout
	_ [unsafe.Sizeof(*new(C.GError))]byte
}

func (err *Error) c() *C.GError {
	return (*C.GError)(unsafe.Pointer(err))
}

func (err *Error) Error() string {
	return C.GoString(err.c().message)
}

type Bytes struct {
	_ structs.HostLayout
	_ [unsafe.Sizeof(*new(C.GBytes))]byte
}

type OptionGroup struct {
	_ structs.HostLayout
	_ [unsafe.Sizeof(*new(C.GOptionGroup))]byte
}
