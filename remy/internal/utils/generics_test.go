package utils

import (
	"testing"
)

func TestTypeName__DifferPointerFromInterface(t *testing.T) {
	type testInterface interface {
		a() bool
	}

	for _, generifyInterface := range [...]bool{true, false} {
		interfaceTypeResult := TypeName[testInterface](generifyInterface)
		pointerTypeResult := TypeName[*testInterface](generifyInterface)
		if interfaceTypeResult == pointerTypeResult {
			t.Error("type names was the same, when it should be different, because of different pointer")
		}
	}
}

func TestTypeNameByReflect__DifferPointerFromInterface(t *testing.T) {
	type testInterface interface {
		a() bool
	}

	for _, generifyInterface := range [...]bool{true, false} {
		interfaceTypeResult := TypeNameByReflect[testInterface](generifyInterface)
		pointerTypeResult := TypeNameByReflect[*testInterface](generifyInterface)
		if interfaceTypeResult == pointerTypeResult {
			t.Error("type names was the same, when it should be different, because of different pointer")
		}
	}
}
