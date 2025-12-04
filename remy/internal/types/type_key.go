package types

import (
	"hash/fnv"
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
	StrKeyElem     string
)

func (k KeyElem[T]) comparable() {
	// Just a stub function
}

func (k KeyElem[T]) ID() uint64 {
	iface := any(k)
	typePointer := uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&iface)))
	return uint64(typePointer)
}

func (k StrKeyElem) comparable() {
	// Just a stub function
}

func (k StrKeyElem) ID() uint64 {
	h := fnv.New64a()
	_, _ = h.Write([]byte(k))
	return h.Sum64()
}
