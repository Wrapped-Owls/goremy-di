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
	if elementType == nil {
		panic(ErrImpossibleIdentifyType)
	}
	return fmt.Sprintf("%s/%s", elementType.PkgPath(), elementType.Name())
}

func GetType[T any]() reflect.Type {
	var typeT T
	foundType := reflect.TypeOf(typeT)
	if foundType == nil {
		// T is an interface
		foundType = reflect.TypeOf(&typeT)
	}
	return foundType
}

func Default[T any]() T {
	var element T
	return element
}
