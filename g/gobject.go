package g

/*
#cgo pkg-config: gobject-2.0
#include <glib-object.h>

void _g_object_dispose(GObject *obj);
*/
import "C"

import (
	"runtime/cgo"
	"slices"
	"structs"
	"unsafe"

	"deedles.dev/gir/internal/util"
)

type Type[T any] struct {
	_ structs.HostLayout
	_ [unsafe.Sizeof(*new(C.GType))]byte
}

func ToType[T any](c uint64) Type[T] {
	return *(*Type[T])(unsafe.Pointer(&c))
}

func (t *Type[T]) c() C.GType {
	return *(*C.GType)(unsafe.Pointer(t))
}

func (t Type[T]) New(props ...any) *T {
	names := make([]*C.char, 0, len(props)/2)
	vals := make([]C.GValue, 0, len(props)/2)
	for n, v := range util.Pairs(slices.Values(props)) {
		names = append(names, C.CString(n.(string)))
		vals = append(vals, v.(C.GValue))
	}

	return (*T)(unsafe.Pointer(C.g_object_new_with_properties(
		t.c(),
		C.guint(len(names)),
		(**C.char)(unsafe.SliceData(names)),
		(*C.GValue)(unsafe.SliceData(vals)),
	)))
}

func (t Type[T]) Cast(obj *TypeInstance) *T {
	v, ok := t.Check(obj)
	if !ok {
		panic("type is not convertible")
	}
	return v
}

func (t Type[T]) Check(obj *TypeInstance) (*T, bool) {
	target := C.g_type_from_name(C.g_type_name_from_instance(obj.c()))
	if C.g_type_is_a(t.c(), target) == 0 && C.g_type_is_a(target, t.c()) == 0 {
		return nil, false
	}
	return (*T)(unsafe.Pointer(obj)), true
}

type TypeClass struct {
	_ structs.HostLayout
	_ [unsafe.Sizeof(*new(C.GTypeClass))]byte
}

func (tc *TypeClass) c() *C.GTypeClass {
	return (*C.GTypeClass)(unsafe.Pointer(tc))
}

func (tc *TypeClass) AsGTypeClass() *TypeClass { return tc }

func (tc *TypeClass) TypeName() string {
	return C.GoString(C.g_type_name_from_class(tc.c()))
}

type TypeInstance struct {
	_ structs.HostLayout
	_ [unsafe.Sizeof(*new(C.GTypeInstance))]byte
}

func (ti *TypeInstance) c() *C.GTypeInstance {
	return (*C.GTypeInstance)(unsafe.Pointer(ti))
}

func (ti *TypeInstance) AsGTypeInstance() *TypeInstance { return ti }

func (ti *TypeInstance) TypeName() string {
	return C.GoString(C.g_type_name_from_instance(ti.c()))
}

type ObjectClass struct {
	_ structs.HostLayout
	TypeClass
	_ [unsafe.Sizeof(*new(C.GObjectClass)) - unsafe.Sizeof(*new(C.GTypeClass))]byte
}

func (class *ObjectClass) c() *C.GObjectClass {
	return (*C.GObjectClass)(unsafe.Pointer(class))
}

func (class *ObjectClass) AsGObjectClass() *ObjectClass { return class }

var _g_object_dispose_quark C.GQuark = C.g_quark_from_static_string(C.CString("_g_object_dispose"))

//export _g_object_dispose
func _g_object_dispose(obj *C.GObject) {
	t := C.g_type_from_name(C.g_type_name_from_instance((*C.GTypeInstance)(unsafe.Pointer(obj))))
	f := cgo.Handle(C.g_type_get_qdata(t, _g_object_dispose_quark)).Value().(func(*Object))
	f((*Object)(unsafe.Pointer(obj)))
}

func (class *ObjectClass) SetDispose(dispose func(*Object)) {
	t := C.g_type_from_name(C.g_type_name_from_class(class.AsGTypeClass().c()))
	h := cgo.Handle(C.g_type_get_qdata(t, _g_object_dispose_quark))
	if h != 0 {
		h.Delete()
	}

	C.g_type_set_qdata(t, _g_object_dispose_quark, C.gpointer(cgo.NewHandle(dispose)))
	class.c().dispose = (*[0]byte)(C._g_object_dispose)
}

type Object struct {
	_ structs.HostLayout
	TypeInstance
	_ [unsafe.Sizeof(*new(C.GObject)) - unsafe.Sizeof(*new(C.GTypeInstance))]byte
}

func (obj *Object) c() *C.GObject {
	return (*C.GObject)(unsafe.Pointer(obj))
}

func (obj *Object) AsGObject() *Object { return obj }

func (obj *Object) Ref() {
	C.g_object_ref(C.gpointer(obj.c()))
}

func (obj *Object) Unref() {
	C.g_object_unref(C.gpointer(obj.c()))
}
