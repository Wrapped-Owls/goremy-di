package utils

import (
	"fmt"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"reflect"
	"strings"
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
	} else {
		value = Default[T]()
	}
	name = fmt.Sprintf("%T", value)
	if name == nilStr {
		name = fmt.Sprintf("%T", &value)
	}
	return
}

// TypeNameByReflect returns a string that defines the name of the given generic type.
func TypeNameByReflect(generifyInterface bool, elementType reflect.Type, isInterface bool) string {
	if elementType == nil {
		panic(ErrImpossibleIdentifyType)
	}

	if generifyInterface && isInterface {
		var builder strings.Builder
		builder.WriteString("interface { ")
		for i := 0; i < elementType.NumMethod(); i++ {
			if i > 0 {
				builder.WriteString("; ")
			}
			builder.WriteString(
				fmt.Sprintf(
					"%s %s", elementType.Method(i).Name, elementType.Method(i).Type,
				),
			)
		}
		builder.WriteString(" }")
		return builder.String()
	}
	return fmt.Sprintf("%s/%s{###}%s", elementType.PkgPath(), elementType.Name(), fmt.Sprint(elementType))
}

func GetElemType(element any) (foundType reflect.Type, isInterface bool) {
	foundType = reflect.TypeOf(element)
	if foundType == nil {
		panic(ErrImpossibleIdentifyType)
	}
	return
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
