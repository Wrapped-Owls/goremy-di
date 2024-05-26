package types

import (
	"reflect"
	"testing"
)

func TestBindKey_AsMapKey(t *testing.T) {
	dependenciesBool := BindDependencies[any]{}

	testList := []struct {
		key   BindKey
		value any
	}{
		{
			key:   KeyElem[bool]{},
			value: true,
		},
		{
			key:   KeyElem[uint8]{},
			value: 42,
		},
		{
			key:   StrKeyElem("key"),
			value: "answers everywhere",
		},
		{
			key:   StrKeyElem("Gomu Gomu no Mi"),
			value: "Hito Hito no Mi, Model: Nika",
		},
		{
			key:   KeyElem[map[string]bool]{},
			value: map[string]bool{"test": true},
		},
		{
			key:   KeyElem[float64]{},
			value: 5216848912325.9924187442,
		},
		{
			key:   KeyElem[string]{},
			value: "this is a little different one",
		},
	}

	for _, subject := range testList {
		dependenciesBool[subject.key] = subject.value
	}

	for _, subject := range testList {
		if val, exists := dependenciesBool[subject.key]; !exists ||
			!reflect.DeepEqual(val, subject.value) {
			t.Errorf(
				"Expected BindDependencies[any] to store and retrieve value correctly, got %v",
				val,
			)
		}
	}
}

// TestComparableBindKey is a test just to improve coverage over the type_key file
func TestComparableBindKey(t *testing.T) {
	testList := []BindKey{KeyElem[bool]{}, StrKeyElem("any")}
	for _, subject := range testList {
		subject.comparable()
	}
}
