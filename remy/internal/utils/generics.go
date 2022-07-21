package utils

import (
	"fmt"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

// nilStr is the default const string representation of a nil type
var nilStr = fmt.Sprintf("%T", nil)

func GetKey[T any](generifyInterface bool) types.BindKey {
	if !generifyInterface {
		return TypeName[T]()
	}
	elementType, isInterface := GetType[T]()
	return TypeNameByReflect(generifyInterface, elementType, isInterface)
}

func GetElemKey(element any, generifyInterface bool) types.BindKey {
	if !generifyInterface {
		return TypeName(element)
	}
	elementType, isInterface := GetElemType(element)
	return TypeNameByReflect(generifyInterface, elementType, isInterface)
}

func TypeName[T any](elements ...T) (name string) {
	var value T
	if len(elements) > 0 {
		value = elements[0]
	}

	name = fmt.Sprintf("%T", value)
	if name == nilStr {
		name = fmt.Sprintf("%T", &value)
	}
	return
}

func Default[T any]() T {
	var element T
	return element
}
