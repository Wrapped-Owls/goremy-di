package utils

import (
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

func shouldGenerify(options injopts.KeyGenOption) bool {
	return options&injopts.KeyOptGenerifyInterface == injopts.KeyOptGenerifyInterface
}

func shouldUseReflection(options injopts.KeyGenOption) bool {
	return options&injopts.KeyOptUseReflectionType == injopts.KeyOptUseReflectionType
}

func shouldPrefixPointer(options injopts.KeyGenOption) bool {
	return options&injopts.KeyOptIgnorePointer != injopts.KeyOptIgnorePointer
}

func GetKey[T any](options injopts.KeyGenOption) types.BindKey {
	generifyInterface := shouldGenerify(options)
	if shouldUseReflection(options) || generifyInterface {
		keyVal, _ := TypeNameByReflection[T](generifyInterface, shouldPrefixPointer(options))
		return types.StrKeyElem(keyVal)
	}

	return NewKeyElem[T]()
}

// IsInterface reports whether T is an interface type.
func IsInterface[T any]() bool {
	var zero T

	// Converting zero to `any` gives:
	//   - true nil  → (nil, nil)        → only happens when T is an interface
	//   - non-nil   → (type != nil, nil) → happens for pointer/slice/map/chan/func
	//
	// Even nil pointers like (*int)(nil) become ( *int, nil ) and are NOT equal to nil.
	// Therefore, only interface types produce a nil interface value.
	if any(zero) != nil {
		return false // concrete type (non-pointer, non-interface)
	}

	return true // must be an interface
}

func NewKeyElem[T any]() types.KeyElem[T] {
	return types.KeyElem[T]{}
}

func GetElemKey(element any, options injopts.KeyGenOption) (types.BindKey, error) {
	if !shouldUseReflection(options) {
		return types.StrKeyElem(""), remyErrs.ErrGetElementTypeRequiresReflectionEnabled
	}

	generifyInterface := shouldGenerify(options)
	keyVal, err := TypeNameByReflection(generifyInterface, shouldPrefixPointer(options), element)
	return types.StrKeyElem(keyVal), err
}
