package gi

//go:generate go tool girgen -o girepository.gen.go girepository.gen

/*
#cgo pkg-config: girepository-2.0

#include <girepository/girepository.h>
*/
import "C"

import (
	"iter"
	"structs"
	"unsafe"

	"deedles.dev/gir/g"
)

var TypeRepository = g.ToType[Repository](uint64(C.gi_repository_get_type()))

type Repository struct {
	_ structs.HostLayout
	g.Object
	//_ [unsafe.Sizeof(*new(C.GIRepository)) - unsafe.Sizeof(*new(C.GObject))]byte
}

func RepositoryNew() *Repository {
	return (*Repository)(unsafe.Pointer(C.gi_repository_new()))
}

func (r *Repository) c() *C.GIRepository {
	return (*C.GIRepository)(unsafe.Pointer(r))
}

type RepositoryLoadFlags int

const (
	RepositoryLoadFlagNone RepositoryLoadFlags = iota
	RepositoryLoadFlagLazy
)

func (r *Repository) Require(namespace, version string, flags RepositoryLoadFlags) (*Typelib, error) {
	cnamespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(cnamespace))

	cversion := C.CString(version)
	defer C.free(unsafe.Pointer(cversion))

	var gerr *C.GError
	tl := (*Typelib)(unsafe.Pointer(C.gi_repository_require(
		r.c(),
		cnamespace,
		cversion,
		C.GIRepositoryLoadFlags(flags),
		&gerr,
	)))
	if gerr != nil {
		return tl, (*g.Error)(unsafe.Pointer(gerr))
	}
	return tl, nil
}

func (r *Repository) GetNInfos(namespace string) uint {
	cnamespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(cnamespace))

	return uint(C.gi_repository_get_n_infos(r.c(), cnamespace))
}

func (r *Repository) GetInfo(namespace string, index uint) *BaseInfo {
	cnamespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(cnamespace))

	return (*BaseInfo)(unsafe.Pointer(C.gi_repository_get_info(r.c(), cnamespace, C.uint(index))))
}

func (r *Repository) GetInfos(namespace string) iter.Seq[*BaseInfo] {
	return func(yield func(*BaseInfo) bool) {
		n := r.GetNInfos(namespace)
		for i := range n {
			info := r.GetInfo(namespace, i)
			if !yield(info) {
				return
			}
		}
	}
}

func (r *Repository) GetCPrefix(namespace string) string {
	cnamespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(cnamespace))

	return C.GoString(C.gi_repository_get_c_prefix(r.c(), cnamespace))
}

var TypeBaseInfo = g.ToType[BaseInfo](uint64(C.gi_base_info_get_type()))

type BaseInfo struct {
	_ structs.HostLayout
	g.TypeInstance
	//_ [unsafe.Sizeof(*new(C.GIBaseInfo)) - unsafe.Sizeof(*new(C.GTypeInstance))]byte
}

func (info *BaseInfo) c() *C.GIBaseInfo {
	return (*C.GIBaseInfo)(unsafe.Pointer(info))
}

func (info *BaseInfo) AsGIBaseInfo() *BaseInfo { return info }

func (info *BaseInfo) GetName() string {
	return C.GoString(C.gi_base_info_get_name(info.c()))
}

func (info *BaseInfo) IterateAttributes() iter.Seq2[string, string] {
	return func(yield func(string, string) bool) {
		var iterator AttributeIter
		for {
			var name, val *C.char
			ok := C.gi_base_info_iterate_attributes(info.c(), iterator.c(), &name, &val)
			if ok == 0 || !yield(C.GoString(name), C.GoString(val)) {
				return
			}
		}
	}
}

func (info *BaseInfo) Ref() {
	C.gi_base_info_ref(unsafe.Pointer(info.c()))
}

func (info *BaseInfo) Unref() {
	C.gi_base_info_unref(unsafe.Pointer(info.c()))
}

var TypeCallableInfo = g.ToType[CallableInfo](uint64(C.gi_callable_info_get_type()))

type CallableInfo struct {
	_ structs.HostLayout
	BaseInfo
	_ [unsafe.Sizeof(*new(C.GICallableInfo)) - unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

func (info *CallableInfo) c() *C.GICallableInfo {
	return (*C.GICallableInfo)(unsafe.Pointer(info))
}

func (info *CallableInfo) AsGICallableInfo() *CallableInfo { return info }

func (info *CallableInfo) IsMethod() bool {
	return C.gi_callable_info_is_method(info.c()) != 0
}

func (info *CallableInfo) IsAsync() bool {
	return C.gi_callable_info_is_async(info.c()) != 0
}

func (info *CallableInfo) GetNArgs() uint {
	return uint(C.gi_callable_info_get_n_args(info.c()))
}

func (info *CallableInfo) GetArg(index uint) *ArgInfo {
	return (*ArgInfo)(unsafe.Pointer(C.gi_callable_info_get_arg(info.c(), C.uint(index))))
}

func (info *CallableInfo) GetArgs() iter.Seq[*ArgInfo] {
	return func(yield func(*ArgInfo) bool) {
		n := info.GetNArgs()
		for i := range n {
			if !yield(info.GetArg(i)) {
				return
			}
		}
	}
}

var TypeFunctionInfo = g.ToType[FunctionInfo](uint64(C.gi_function_info_get_type()))

type FunctionInfo struct {
	_ structs.HostLayout
	CallableInfo
	_ [unsafe.Sizeof(*new(C.GIFunctionInfo)) - unsafe.Sizeof(*new(C.GICallableInfo))]byte
}

func (info *FunctionInfo) c() *C.GIFunctionInfo {
	return (*C.GIFunctionInfo)(unsafe.Pointer(info))
}

func (info *FunctionInfo) AsGIFunctionInfo() *FunctionInfo { return info }

type FunctionInfoFlags int

const (
	FunctionInfoFlagsNone FunctionInfoFlags = 0
	FunctionIsMethod      FunctionInfoFlags = 1 << (iota - 1)
	FunctionIsConstructor
	FunctionIsGetter
	FunctionIsSetter
	FunctionWrapsVFunc
	FunctionIsAsync
)

func (info *FunctionInfo) GetFlags() FunctionInfoFlags {
	return FunctionInfoFlags(C.gi_function_info_get_flags(info.c()))
}

var TypeConstantInfo = g.ToType[ConstantInfo](uint64(C.gi_constant_info_get_type()))

type ConstantInfo struct {
	_ structs.HostLayout
	BaseInfo
	_ [unsafe.Sizeof(*new(C.GIConstantInfo)) - unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

func (info *ConstantInfo) c() *C.GIConstantInfo {
	return (*C.GIConstantInfo)(unsafe.Pointer(info))
}

func (info *ConstantInfo) AsGIConstantInfo() *ConstantInfo { return info }

var TypeRegisteredTypeInfo = g.ToType[RegisteredTypeInfo](uint64(C.gi_registered_type_info_get_type()))

type RegisteredTypeInfo struct {
	_ structs.HostLayout
	BaseInfo
	_ [unsafe.Sizeof(*new(C.GIRegisteredTypeInfo)) - unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

func (info *RegisteredTypeInfo) c() *C.GIRegisteredTypeInfo {
	return (*C.GIRegisteredTypeInfo)(unsafe.Pointer(info))
}

func (info *RegisteredTypeInfo) AsGIRegisteredTypeInfo() *RegisteredTypeInfo { return info }

func (info *RegisteredTypeInfo) GetTypeName() string {
	return C.GoString(C.gi_registered_type_info_get_type_name(info.c()))
}

var TypeObjectInfo = g.ToType[ObjectInfo](uint64(C.gi_object_info_get_type()))

type ObjectInfo struct {
	_ structs.HostLayout
	BaseInfo
	_ [unsafe.Sizeof(*new(C.GIObjectInfo)) - unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

func (info *ObjectInfo) c() *C.GIObjectInfo {
	return (*C.GIObjectInfo)(unsafe.Pointer(info))
}

func (info *ObjectInfo) AsGIObjectInfo() *ObjectInfo { return info }

func (info *ObjectInfo) GetNMethods() uint {
	return uint(C.gi_object_info_get_n_methods(info.c()))
}

func (info *ObjectInfo) GetMethod(index uint) *FunctionInfo {
	return (*FunctionInfo)(unsafe.Pointer(C.gi_object_info_get_method(info.c(), C.uint(index))))
}

func (info *ObjectInfo) GetMethods() iter.Seq[*FunctionInfo] {
	return func(yield func(*FunctionInfo) bool) {
		n := info.GetNMethods()
		for i := range n {
			if !yield(info.GetMethod(i)) {
				return
			}
		}
	}
}

var TypeStructInfo = g.ToType[StructInfo](uint64(C.gi_struct_info_get_type()))

type StructInfo struct {
	_ structs.HostLayout
	RegisteredTypeInfo
	_ [unsafe.Sizeof(*new(C.GIStructInfo)) - unsafe.Sizeof(*new(C.GIRegisteredTypeInfo))]byte
}

func (info *StructInfo) c() *C.GIStructInfo {
	return (*C.GIStructInfo)(unsafe.Pointer(info))
}

func (info *StructInfo) AsGIStructInfo() *StructInfo { return info }

func (info *StructInfo) GetNMethods() uint {
	return uint(C.gi_struct_info_get_n_methods(info.c()))
}

func (info *StructInfo) GetMethod(index uint) *FunctionInfo {
	return (*FunctionInfo)(unsafe.Pointer(C.gi_struct_info_get_method(info.c(), C.uint(index))))
}

func (info *StructInfo) GetMethods() iter.Seq[*FunctionInfo] {
	return func(yield func(*FunctionInfo) bool) {
		n := info.GetNMethods()
		for i := range n {
			if !yield(info.GetMethod(i)) {
				return
			}
		}
	}
}

func (info *StructInfo) GetNFields() uint {
	return uint(C.gi_struct_info_get_n_fields(info.c()))
}

func (info *StructInfo) GetField(index uint) *FieldInfo {
	return (*FieldInfo)(unsafe.Pointer(C.gi_struct_info_get_field(info.c(), C.uint(index))))
}

func (info *StructInfo) GetFields() iter.Seq[*FieldInfo] {
	return func(yield func(*FieldInfo) bool) {
		n := info.GetNFields()
		for i := range n {
			if !yield(info.GetField(i)) {
				return
			}
		}
	}
}

func (info *StructInfo) GetSize() uint64 {
	return uint64(C.gi_struct_info_get_size(info.c()))
}

func (info *StructInfo) IsForeign() bool {
	return C.gi_struct_info_is_foreign(info.c()) != 0
}

type FieldInfo struct {
	_ structs.HostLayout
	BaseInfo
	_ [unsafe.Sizeof(*new(C.GIFieldInfo)) - unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

func (info *FieldInfo) c() *C.GIFieldInfo {
	return (*C.GIFieldInfo)(unsafe.Pointer(info))
}

type ArgInfo struct {
	_ structs.HostLayout
	BaseInfo
	_ [unsafe.Sizeof(*new(C.GIArgInfo)) - unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

func (info *ArgInfo) c() *C.GIArgInfo {
	return (*C.GIArgInfo)(unsafe.Pointer(info))
}

func (info *ArgInfo) AsGIArgInfo() *ArgInfo { return info }
