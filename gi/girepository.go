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
)

func (r *Repository) GetInfos(namespace string) iter.Seq2[uint32, *BaseInfo] {
	return func(yield func(uint32, *BaseInfo) bool) {
		n := r.GetNInfos(namespace)
		for i := range n {
			info := r.GetInfo(namespace, i)
			if !yield(i, info) {
				return
			}
		}
	}
}

func (info *CallableInfo) GetArgs() iter.Seq2[uint32, *ArgInfo] {
	return func(yield func(uint32, *ArgInfo) bool) {
		n := info.GetNArgs()
		for i := range n {
			if !yield(i, info.GetArg(i)) {
				return
			}
		}
	}
}

func (info *ObjectInfo) GetMethods() iter.Seq2[uint32, *FunctionInfo] {
	return func(yield func(uint32, *FunctionInfo) bool) {
		n := info.GetNMethods()
		for i := range n {
			if !yield(i, info.GetMethod(i)) {
				return
			}
		}
	}
}

func (info *StructInfo) GetMethods() iter.Seq2[uint32, *FunctionInfo] {
	return func(yield func(uint32, *FunctionInfo) bool) {
		n := info.GetNMethods()
		for i := range n {
			if !yield(i, info.GetMethod(i)) {
				return
			}
		}
	}
}

func (info *StructInfo) GetFields() iter.Seq2[uint32, *FieldInfo] {
	return func(yield func(uint32, *FieldInfo) bool) {
		n := info.GetNFields()
		for i := range n {
			if !yield(i, info.GetField(i)) {
				return
			}
		}
	}
}

func (info *EnumInfo) GetValues() iter.Seq2[uint32, *ValueInfo] {
	return func(yield func(uint32, *ValueInfo) bool) {
		n := info.GetNValues()
		for i := range n {
			v := info.GetValue(i)
			if !yield(i, v) {
				return
			}
		}
	}
}

type Argument struct {
	_ structs.HostLayout
	_ [unsafe.Sizeof(*new(C.GIArgument))]byte
}
