package utils

import (
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
		return types.StrKeyElem(
			TypeNameByReflection[T](generifyInterface, shouldPrefixPointer(options)),
		)
	}

	return types.KeyElem[T]{}
}

func GetElemKey(element any, options injopts.KeyGenOption) types.BindKey {
	generifyInterface := shouldGenerify(options)
	return types.StrKeyElem(
		TypeNameByReflection(generifyInterface, shouldPrefixPointer(options), element),
	)
}
