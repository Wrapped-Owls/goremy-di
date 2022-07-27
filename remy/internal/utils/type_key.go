package utils

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

func GetKey[T any](options types.ReflectionOptions) types.BindKey {
	if options.UseReflectionType {
		return types.BindKey(TypeNameByReflect[T](options.GenerifyInterface))
	}
	return types.BindKey(TypeName[T](options.GenerifyInterface))
}

func GetElemKey(element any, options types.ReflectionOptions) types.BindKey {
	if options.UseReflectionType {
		return types.BindKey(TypeNameByReflect(options.GenerifyInterface, element))
	}
	return types.BindKey(TypeName(options.GenerifyInterface, element))
}
