package types

import (
	"unsafe"
)

type (
	BindKey interface {
		comparable()
		ID() uint64
	}

	BindDependencies[T any] map[BindKey]T
	DependencyGraph         struct {
		UnnamedDependency BindDependencies[bool]
		NamedDependency   BindDependencies[map[string]bool]
	}

	KeyElem[T any] struct{}
)

func (k KeyElem[T]) comparable() {
	// Just a stub function
}

func (k KeyElem[T]) ID() uint64 {
	iface := any(k)
	typePointer := uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&iface)))
	return uint64(typePointer)
}
