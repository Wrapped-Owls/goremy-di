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
func TypeNameByReflect[T any](generifyInterface, identifyPointer bool, elements ...T) string {
	var (
		elementType reflect.Type
		isInterface bool
		value       T
	)
	if len(elements) > 0 {
		value = elements[0]
		elementType, isInterface = GetElemType(value)
	} else {
		elementType, isInterface = GetType[T]()
	}
	if elementType == nil {
		panic(ErrImpossibleIdentifyType)
	}

	if generifyInterface && isInterface {
		return buildDuckInterfaceType(elementType)
	}

	name := fmt.Sprintf(
		"%s/%s{###}%s",
		elementType.PkgPath(), elementType.Name(), fmt.Sprint(elementType),
	)
	if !isInterface && identifyPointer {
		// check if is a pointer
		if reflect.ValueOf(value).Kind() == reflect.Ptr {
			name = fmt.Sprintf("pointer_&%s", name)
		}
	}
	return name
}

func GetElemType[T any](element T) (foundType reflect.Type, isInterface bool) {
	foundType = reflect.TypeOf(element)
	if foundType == nil {
		// element is an interface
		isInterface = true
		foundType = reflect.TypeOf(&element)
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
