package utils

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/keyopts"
)

func shouldGenerify(options keyopts.GenOption) bool {
	return options&keyopts.KeyOptGenerifyInterface == keyopts.KeyOptGenerifyInterface
}

func shouldUseReflection(options keyopts.GenOption) bool {
	return options&keyopts.KeyOptUseReflectionType == keyopts.KeyOptUseReflectionType
}

func shouldPrefixPointer(options keyopts.GenOption) bool {
	return options&keyopts.KeyOptIgnorePointer != keyopts.KeyOptIgnorePointer
}

func GetKey[T any](options keyopts.GenOption) types.BindKey {
	if shouldUseReflection(options) {
		return types.BindKey(TypeNameByReflect[T](shouldGenerify(options), shouldPrefixPointer(options)))
	}
	return types.BindKey(TypeName[T](shouldGenerify(options), shouldPrefixPointer(options)))
}

func GetElemKey(element any, options keyopts.GenOption) types.BindKey {
	if shouldUseReflection(options) {
		return types.BindKey(TypeNameByReflect(shouldGenerify(options), shouldPrefixPointer(options), element))
	}
	return types.BindKey(TypeName(shouldGenerify(options), shouldPrefixPointer(options), element))
}
