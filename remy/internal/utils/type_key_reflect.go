//go:build reflect

package utils

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

func GetKey[T any](generifyInterface bool) types.BindKey {
	elementType, isInterface := GetType[T]()
	return TypeNameByReflect(generifyInterface, elementType, isInterface)
}

func GetElemKey(element any, generifyInterface bool) types.BindKey {
	elementType, isInterface := GetElemType(element)
	return TypeNameByReflect(generifyInterface, elementType, isInterface)
}
