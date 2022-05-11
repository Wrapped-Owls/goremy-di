package utils

import (
	"fmt"
	"github.com/wrapped-owls/fitpiece/gotalaria/internal/types"
	"reflect"
)

func GetKey[T any]() types.BindKey {
	return TypeName[T]()
}

// TypeName returns a string that defines the name of the given generic type.
//
// TODO: Create a typeNameInterface that generates the name based on interface methods signature,
// TODO: so it can be used without importing interfaces (add a flag for it)
func TypeName[T any]() string {
	elementType, _ := GetType[T]()
	if elementType == nil {
		panic(ErrImpossibleIdentifyType)
	}
	return fmt.Sprintf("%s/%s{###}%s", elementType.PkgPath(), elementType.Name(), fmt.Sprint(elementType))
}

func GetType[T any]() (foundType reflect.Type, isInterface bool) {
	var typeT T
	foundType = reflect.TypeOf(typeT)
	if foundType == nil {
		// T is an interface
		isInterface = true
		foundType = reflect.TypeOf(&typeT)
	}
	return
}

func Default[T any]() T {
	var element T
	return element
}
