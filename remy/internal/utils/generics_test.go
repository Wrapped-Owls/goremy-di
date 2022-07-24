package utils

import (
	"testing"
)

func TestTypeName__DifferPointerFromInterface(t *testing.T) {
	type testInterface interface {
		a() bool
	}

	interfaceTypeResult := TypeName[testInterface](false)
	pointerTypeResult := TypeName[*testInterface](false)
	if interfaceTypeResult == pointerTypeResult {
		t.Error("type names was the same, when it should be different, because of different pointer")
	}
}

func TestTypeNameByReflect__DifferPointerFromInterface(t *testing.T) {
	type testInterface interface {
		a() bool
	}

	interfaceTypeResult := TypeNameByReflect[testInterface](false)
	pointerTypeResult := TypeNameByReflect[*testInterface](false)
	if interfaceTypeResult == pointerTypeResult {
		t.Error("type names was the same, when it should be different, because of different pointer")
	}
}
