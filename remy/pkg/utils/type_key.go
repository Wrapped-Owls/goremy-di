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
	if shouldUseReflection(options) {
		return types.BindKey(
			TypeNameByReflect[T](shouldGenerify(options), shouldPrefixPointer(options)),
		)
	}
	return types.BindKey(TypeName[T](shouldGenerify(options), shouldPrefixPointer(options)))
}

func GetElemKey(element any, options injopts.KeyGenOption) types.BindKey {
	if shouldUseReflection(options) {
		return types.BindKey(
			TypeNameByReflect(shouldGenerify(options), shouldPrefixPointer(options), element),
		)
	}
	return types.BindKey(TypeName(shouldGenerify(options), shouldPrefixPointer(options), element))
}
