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

func (info *BaseInfo) GetNamespace() string {
	return C.GoString(C.gi_base_info_get_namespace(info.c()))
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

func (info *RegisteredTypeInfo) GetGType() g.Type[g.TypeInstance] {
	return g.ToType[g.TypeInstance](uint64(C.gi_registered_type_info_get_g_type(info.c())))
}

var TypeObjectInfo = g.ToType[ObjectInfo](uint64(C.gi_object_info_get_type()))

type ObjectInfo struct {
	_ structs.HostLayout
	RegisteredTypeInfo
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

type FieldInfo struct {
	_ structs.HostLayout
	BaseInfo
	_ [unsafe.Sizeof(*new(C.GIFieldInfo)) - unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

func (info *FieldInfo) c() *C.GIFieldInfo {
	return (*C.GIFieldInfo)(unsafe.Pointer(info))
}

var TypeArgInfo = g.ToType[ArgInfo](uint64(C.gi_arg_info_get_type()))

type ArgInfo struct {
	_ structs.HostLayout
	BaseInfo
	_ [unsafe.Sizeof(*new(C.GIArgInfo)) - unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

func (info *ArgInfo) c() *C.GIArgInfo {
	return (*C.GIArgInfo)(unsafe.Pointer(info))
}

func (info *ArgInfo) AsGIArgInfo() *ArgInfo { return info }

func (info *ArgInfo) GetTypeInfo() *TypeInfo {
	return (*TypeInfo)(unsafe.Pointer(C.gi_arg_info_get_type_info(info.c())))
}

func (info *ArgInfo) IsReturnValue() bool {
	return C.gi_arg_info_is_return_value(info.c()) != 0
}

type Direction int

const (
	DirectionIn Direction = iota
	DirectionOut
	DirectionInout
)

func (info *ArgInfo) GetDirection() Direction {
	return Direction(C.gi_arg_info_get_direction(info.c()))
}

func (info *ArgInfo) IsSkip() bool {
	return C.gi_arg_info_is_skip(info.c()) != 0
}

type Transfer int

const (
	TransferNothing Transfer = iota
	TransferContainer
	TransferEverything
)

func (info *ArgInfo) GetOwnershipTransfer() Transfer {
	return Transfer(C.gi_arg_info_get_ownership_transfer(info.c()))
}

type TypeInfo struct {
	_ structs.HostLayout
	BaseInfo
	_ [unsafe.Sizeof(*new(C.GITypeInfo)) - unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

func (info *TypeInfo) c() *C.GITypeInfo {
	return (*C.GITypeInfo)(unsafe.Pointer(info))
}

func (info *TypeInfo) AsGITypeInfo() *TypeInfo {
	return info
}

type TypeTag int

const (
	TypeTagVoid TypeTag = iota
	TypeTagBoolean
	TypeTagInt8
	TypeTagUint8
	TypeTagInt16
	TypeTagUint16
	TypeTagInt32
	TypeTagUint32
	TypeTagInt64
	TypeTagUint64
	TypeTagFloat
	TypeTagDouble
	TypeTagGtype
	TypeTagUtf
	TypeTagFilename
	TypeTagArray
	TypeTagInterface
	TypeTagGlist
	TypeTagGslist
	TypeTagGhash
	TypeTagError
	TypeTagUnichar
)

func (tag TypeTag) String() string {
	return [...]string{
		"Void",
		"Boolean",
		"Int",
		"Uint",
		"Int",
		"Uint",
		"Int",
		"Uint",
		"Int",
		"Uint",
		"Float",
		"Double",
		"Gtype",
		"Utf",
		"Filename",
		"Array",
		"Interface",
		"Glist",
		"Gslist",
		"Ghash",
		"Error",
		"Unichar",
	}[tag]
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
