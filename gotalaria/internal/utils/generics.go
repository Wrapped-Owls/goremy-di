package utils

import (
	"fmt"
	"reflect"
)

func TypeName[T any]() string {
	var typeT T
	elementType := reflect.TypeOf(typeT)
	return fmt.Sprintf("%s/%s", elementType.PkgPath(), elementType.Name())
}
