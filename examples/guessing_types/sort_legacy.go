//go:build !go1.21

package main

import (
	"reflect"
	"sort"
	"strings"
)

func sortByReflection[T any](listOfCalculators []T) {
	sort.Slice(listOfCalculators, func(i, j int) bool {
		aType := reflect.TypeOf(listOfCalculators[i])
		bType := reflect.TypeOf(listOfCalculators[j])
		return strings.Compare(aType.Name(), bType.Name()) < 0
	})
}
