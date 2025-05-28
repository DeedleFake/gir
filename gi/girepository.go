package gi

/*
#cgo pkg-config: girepository-2.0

#include <girepository/girepository.h>
*/
import "C"
import (
	"errors"
	"iter"
	"unsafe"
)

type Repository struct {
	_ [unsafe.Sizeof(*new(C.GIRepository))]byte
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
		return tl, errors.New("test")
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

type BaseInfo struct {
	_ [unsafe.Sizeof(*new(C.GIBaseInfo))]byte
}

type Typelib struct {
	_ [unsafe.Sizeof(*new(C.GITypelib))]byte
}
