package utils

import (
	"fmt"
	"reflect"
	"strings"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
)

func buildDuckInterfaceType(elementType reflect.Type) string {
	var builder strings.Builder
	builder.WriteString("interface { ")
	if elementType.Kind() == reflect.Ptr {
		if newElemType := elementType.Elem(); newElemType != nil {
			elementType = newElemType
		}
	}

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

// TypeNameByReflection returns a string that defines the name of the given generic type.
func TypeNameByReflection[T any](
	generifyInterface, identifyPointer bool, elements ...T,
) (string, error) {
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
		return "", &remyErrs.ErrImpossibleIdentifyType{Type: (*T)(nil)}
	}

	if isInterface && generifyInterface {
		return buildDuckInterfaceType(elementType), nil
	}

	name := elementType.PkgPath() + "/" + elementType.Name() + "{###}" + elementType.String()
	if !isInterface && identifyPointer {
		// check if is a pointer
		if reflect.ValueOf(value).Kind() == reflect.Ptr {
			name = "pointer_&" + name
		}
	}
	return name, nil
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
	return GetElemType(typeT)
}
