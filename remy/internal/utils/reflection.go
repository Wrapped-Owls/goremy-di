package utils

import (
	"fmt"
	"reflect"
	"strings"
)

func buildDuckInterfaceType(elementType reflect.Type) string {
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

// TypeNameByReflect returns a string that defines the name of the given generic type.
func TypeNameByReflect(generifyInterface bool, elementType reflect.Type, isInterface bool) string {
	if elementType == nil {
		panic(ErrImpossibleIdentifyType)
	}

	if generifyInterface && isInterface {
		return buildDuckInterfaceType(elementType)
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
