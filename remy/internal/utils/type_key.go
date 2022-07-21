//go:build !reflect

package utils

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

func GetKey[T any](generifyInterface bool) types.BindKey {
	return TypeName[T](generifyInterface)
}

func GetElemKey(element any, generifyInterface bool) types.BindKey {
	return TypeName(generifyInterface, element)
}
