package utils

import (
	"fmt"
	"github.com/wrapped-owls/fitpiece/gotalaria/internal/types"
	"reflect"
)

func GetKey[T any]() types.BindKey {
	return TypeName[T]()
}

func TypeName[T any]() string {
	elementType := GetType[T]()
	return fmt.Sprintf("%s/%s", elementType.PkgPath(), elementType.Name())
}

func GetType[T any]() reflect.Type {
	var typeT T
	return reflect.TypeOf(typeT)
}

func Default[T any]() T {
	var element T
	return element
}
