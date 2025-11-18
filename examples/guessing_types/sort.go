//go:build go1.21

package main

import (
	"reflect"
	"slices"
	"strings"
)

func sortByReflection[T any](listOfCalculators []T) {
	slices.SortFunc(listOfCalculators, func(a, b T) int {
		aType := reflect.TypeOf(a)
		bType := reflect.TypeOf(b)
		return strings.Compare(aType.Name(), bType.Name())
	})
}
