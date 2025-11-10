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
