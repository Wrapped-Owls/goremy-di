package utils

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

func GetKey[T any](options types.ReflectionOptions) types.BindKey {
	if options.UseReflectionType {
		return TypeNameByReflect[T](options.GenerifyInterface)
	}
	return TypeName[T](options.GenerifyInterface)
}

func GetElemKey(element any, options types.ReflectionOptions) types.BindKey {
	if options.UseReflectionType {
		return TypeNameByReflect(options.GenerifyInterface, element)
	}
	return TypeName(options.GenerifyInterface, element)
}
