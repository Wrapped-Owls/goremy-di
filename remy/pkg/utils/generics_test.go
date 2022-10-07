package utils

import "testing"

const typeNameErr = "type names was the same, when it should be different, because of different pointer"

func TestTypeName__DifferPointerFromInterface(t *testing.T) {
	type testInterface interface {
		a() bool
	}

	for _, generifyInterface := range [...]bool{true, false} {
		interfaceTypeResult := TypeName[testInterface](generifyInterface, true)
		pointerTypeResult := TypeName[*testInterface](generifyInterface, true)
		doublePointerTypeResult := TypeName[**testInterface](generifyInterface, true)
		if interfaceTypeResult == pointerTypeResult {
			t.Error(typeNameErr)
		}
		if doublePointerTypeResult == pointerTypeResult {
			t.Error(typeNameErr)
		}
	}
}

func TestTypeNameByReflect__DifferPointerFromInterface(t *testing.T) {
	type testInterface interface {
		a() bool
	}

	for _, generifyInterface := range [...]bool{true, false} {
		interfaceTypeResult := TypeNameByReflect[testInterface](generifyInterface, true)
		pointerTypeResult := TypeNameByReflect[*testInterface](generifyInterface, true)
		doublePointerTypeResult := TypeNameByReflect[**testInterface](generifyInterface, true)
		if interfaceTypeResult == pointerTypeResult {
			t.Error(typeNameErr)
		}
		if doublePointerTypeResult == pointerTypeResult {
			t.Error(typeNameErr)
		}
	}
}
