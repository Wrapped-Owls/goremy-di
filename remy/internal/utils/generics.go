package utils

import (
	"fmt"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"reflect"
	"strings"
)

func GetKey[T any](generifyInterface bool) types.BindKey {
	return TypeName[T](generifyInterface)
}

// TypeName returns a string that defines the name of the given generic type.
func TypeName[T any](generifyInterface bool) string {
	elementType, isInterface := GetType[T]()
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
