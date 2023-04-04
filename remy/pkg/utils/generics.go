package utils

import (
	"fmt"
	"strings"
)

// nilStr is the default const string representation of a nil type
var nilStr = fmt.Sprintf("%T", nil)

func interfaceTypeName[T any](shouldGenerify bool, element T) (name string) {
	if shouldGenerify {
		if elementType, isInterface := GetElemType(element); isInterface {
			return buildDuckInterfaceType(elementType.Elem())
		}
	}
	name = fmt.Sprintf("%T", element)
	if name == nilStr {
		name = fmt.Sprintf("%T", &element)
	}
	return
}

func TypeName[T any](shouldGenerify, identifyPointer bool, elements ...T) (name string) {
	var value T
	if len(elements) > 0 {
		value = elements[0]
	}

	name = fmt.Sprintf("%T", value)
	if name == nilStr {
		name = interfaceTypeName(shouldGenerify, value)
	} else if strings.HasPrefix(name, "*") && identifyPointer {
		name = fmt.Sprintf("pointer_&%s", name)
	}
	return
}
