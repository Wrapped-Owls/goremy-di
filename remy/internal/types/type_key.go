package types

import (
	"unsafe"
)

type (
	BindKey interface {
		ID() uint64
	}

	BindDependencies[T any] map[BindKey]T
	DependencyGraph         struct {
		UnnamedDependency BindDependencies[bool]
		NamedDependency   BindDependencies[map[string]bool]
	}

	KeyElem[T any] struct{}
)

// ID returns a stable, type-unique identifier for KeyElem[T].
// It uses the classic “interface-header” technique:
//
//	Converting k to `any` produces an empty interface value whose
//	first word is a pointer to T’s runtime type descriptor.
//	That pointer is unique per concrete type and constant for the
//	lifetime of the program (within a single binary, no plugins).
//
// We read that first word directly, without reflection, to obtain
// the type identity.
//
// For reference on interface layout, see:
//   - "The Laws of Reflection" (Go Blog):
//     https://go.dev/blog/laws-of-reflection#the-representation-of-an-interface
func (k KeyElem[T]) ID() uint64 {
	iface := any(k) // interface header = [typeptr][dataptr]
	typePointer := uintptr(*(*unsafe.Pointer)(unsafe.Pointer(&iface)))
	return uint64(typePointer)
}
