package utils

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

func GetKey[T any](options types.ReflectionOptions) types.BindKey {
	if options.UseReflectionType {
		elementType, isInterface := GetType[T]()
		return TypeNameByReflect(options.GenerifyInterface, elementType, isInterface)
	}
	return TypeName[T](options.GenerifyInterface)
}

func GetElemKey(element any, options types.ReflectionOptions) types.BindKey {
	if options.UseReflectionType {
		elementType, isInterface := GetElemType(element)
		return TypeNameByReflect(options.GenerifyInterface, elementType, isInterface)
	}
	return TypeName(options.GenerifyInterface, element)
}
