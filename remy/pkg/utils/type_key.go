package utils

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

// IsInterface reports whether T is an interface type.
func IsInterface[T any]() bool {
	var zero T

	// Converting zero to `any` gives:
	//   - true nil  → (nil, nil)        → only happens when T is an interface
	//   - non-nil   → (type != nil, nil) → happens for pointer/slice/map/chan/func
	//
	// Even nil pointers like (*int)(nil) become ( *int, nil ) and are NOT equal to nil.
	// Therefore, only interface types produce a nil interface value.
	//
	// For reference, see:
	//   - Go Blog: “The Laws of Reflection” – explains how interface values internally store a (value, type) pair
	//       https://go.dev/blog/laws-of-reflection#the-representation-of-an-interface
	//   - Go FAQ: “Why typed nil != nil interface” – explains nil vs typed-nil in interfaces
	//       https://go.dev/doc/faq#nil_error
	if any(zero) != nil {
		return false // concrete type (non-pointer, non-interface)
	}

	return true // must be an interface
}

func NewKeyElem[T any]() types.KeyElem[T] {
	return types.KeyElem[T]{}
}
