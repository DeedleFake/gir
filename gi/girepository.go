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

func RepositoryNew() *Repository {
	return (*Repository)(unsafe.Pointer(C.gi_repository_new()))
}

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

func (r *Repository) GetInfos(namespace string) iter.Seq2[uint, *BaseInfo] {
	return func(yield func(uint, *BaseInfo) bool) {
		n := r.GetNInfos(namespace)
		for i := range n {
			info := r.GetInfo(namespace, i)
			if !yield(i, info) {
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

func (info *BaseInfo) GetNamespace() string {
	return C.GoString(C.gi_base_info_get_namespace(info.c()))
}

func (info *BaseInfo) Ref() {
	C.gi_base_info_ref(unsafe.Pointer(info.c()))
}

func (info *BaseInfo) Unref() {
	C.gi_base_info_unref(unsafe.Pointer(info.c()))
}

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

func (info *CallableInfo) GetArgs() iter.Seq2[uint, *ArgInfo] {
	return func(yield func(uint, *ArgInfo) bool) {
		n := info.GetNArgs()
		for i := range n {
			if !yield(i, info.GetArg(i)) {
				return
			}
		}
	}
}

func (info *CallableInfo) GetReturnType() *TypeInfo {
	return (*TypeInfo)(unsafe.Pointer(C.gi_callable_info_get_return_type(info.c())))
}

func (info *CallableInfo) CanThrowGerror() bool {
	return C.gi_callable_info_can_throw_gerror(info.c()) != 0
}

func (info *FunctionInfo) GetFlags() FunctionInfoFlags {
	return FunctionInfoFlags(C.gi_function_info_get_flags(info.c()))
}

func (info *RegisteredTypeInfo) GetTypeName() string {
	return C.GoString(C.gi_registered_type_info_get_type_name(info.c()))
}

func (info *RegisteredTypeInfo) GetGType() g.Type[g.TypeInstance] {
	return g.Type[g.TypeInstance](C.gi_registered_type_info_get_g_type(info.c()))
}

func (info *RegisteredTypeInfo) GetTypeInitFunctionName() string {
	return C.GoString(C.gi_registered_type_info_get_type_init_function_name(info.c()))
}

func (info *ObjectInfo) GetNMethods() uint {
	return uint(C.gi_object_info_get_n_methods(info.c()))
}

func (info *ObjectInfo) GetMethod(index uint) *FunctionInfo {
	return (*FunctionInfo)(unsafe.Pointer(C.gi_object_info_get_method(info.c(), C.uint(index))))
}

func (info *ObjectInfo) GetMethods() iter.Seq2[uint, *FunctionInfo] {
	return func(yield func(uint, *FunctionInfo) bool) {
		n := info.GetNMethods()
		for i := range n {
			if !yield(i, info.GetMethod(i)) {
				return
			}
		}
	}
}

func (info *ObjectInfo) GetParent() *ObjectInfo {
	return (*ObjectInfo)(unsafe.Pointer(C.gi_object_info_get_parent(info.c())))
}

func (info *StructInfo) GetNMethods() uint {
	return uint(C.gi_struct_info_get_n_methods(info.c()))
}

func (info *StructInfo) GetMethod(index uint) *FunctionInfo {
	return (*FunctionInfo)(unsafe.Pointer(C.gi_struct_info_get_method(info.c(), C.uint(index))))
}

func (info *StructInfo) GetMethods() iter.Seq2[uint, *FunctionInfo] {
	return func(yield func(uint, *FunctionInfo) bool) {
		n := info.GetNMethods()
		for i := range n {
			if !yield(i, info.GetMethod(i)) {
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

func (info *StructInfo) GetFields() iter.Seq2[uint, *FieldInfo] {
	return func(yield func(uint, *FieldInfo) bool) {
		n := info.GetNFields()
		for i := range n {
			if !yield(i, info.GetField(i)) {
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

func (info *ArgInfo) GetTypeInfo() *TypeInfo {
	return (*TypeInfo)(unsafe.Pointer(C.gi_arg_info_get_type_info(info.c())))
}

func (info *ArgInfo) IsReturnValue() bool {
	return C.gi_arg_info_is_return_value(info.c()) != 0
}

func (info *ArgInfo) GetDirection() Direction {
	return Direction(C.gi_arg_info_get_direction(info.c()))
}

func (info *ArgInfo) IsSkip() bool {
	return C.gi_arg_info_is_skip(info.c()) != 0
}

func (info *ArgInfo) GetOwnershipTransfer() Transfer {
	return Transfer(C.gi_arg_info_get_ownership_transfer(info.c()))
}

func (info *TypeInfo) GetTag() TypeTag {
	return TypeTag(C.gi_type_info_get_tag(info.c()))
}

func (info *TypeInfo) GetStorageType() TypeTag {
	return TypeTag(C.gi_type_info_get_storage_type(info.c()))
}

func (info *TypeInfo) GetParamType(n uint) *TypeInfo {
	return (*TypeInfo)(unsafe.Pointer(C.gi_type_info_get_param_type(info.c(), C.uint(n))))
}

func (info *TypeInfo) GetInterface() *BaseInfo {
	return (*BaseInfo)(unsafe.Pointer(C.gi_type_info_get_interface(info.c())))
}

func (info *TypeInfo) IsPointer() bool {
	return C.gi_type_info_is_pointer(info.c()) != 0
}

func (info *EnumInfo) GetNValues() uint {
	return uint(C.gi_enum_info_get_n_values(info.c()))
}

func (info *EnumInfo) GetValue(index uint) *ValueInfo {
	return (*ValueInfo)(unsafe.Pointer(C.gi_enum_info_get_value(info.c(), C.uint(index))))
}

func (info *EnumInfo) GetValues() iter.Seq2[uint, *ValueInfo] {
	return func(yield func(uint, *ValueInfo) bool) {
		n := info.GetNValues()
		for i := range n {
			v := info.GetValue(i)
			if !yield(i, v) {
				return
			}
		}
	}
}

func (info *ValueInfo) GetValue() int64 {
	return int64(C.gi_value_info_get_value(info.c()))
}

func (info *TypeInfo) GetArrayType() ArrayType {
	return ArrayType(C.gi_type_info_get_array_type(info.c()))
}

func (info *BaseInfo) GetContainer() *BaseInfo {
	return (*BaseInfo)(unsafe.Pointer(C.gi_base_info_get_container(info.c())))
}

func (info *ObjectInfo) GetRefFunctionName() string {
	return C.GoString(C.gi_object_info_get_ref_function_name(info.c()))
}

type Argument struct {
	_ structs.HostLayout
	_ [unsafe.Sizeof(*new(C.GIArgument))]byte
}

func (info *TypeInfo) GetArrayLengthIndex() (uint, bool) {
	var i C.uint
	ok := C.gi_type_info_get_array_length_index(info.c(), &i)
	return uint(i), ok != 0
}
