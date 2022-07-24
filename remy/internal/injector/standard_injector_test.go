package injector

import (
	"fmt"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

func TestStdInjector_SubInjector(t *testing.T) {
	const strFirstHalf = "the counter is at"
	parent := New(false, types.ReflectionOptions{})
	subInjector := parent.SubInjector(false)

	var counter uint8 = 0
	Register(parent, binds.Factory(func(retriever types.DependencyRetriever) uint8 {
		counter++
		return counter
	}))

	Register(subInjector, binds.Factory(func(retriever types.DependencyRetriever) string {
		return fmt.Sprintf("%s %d", strFirstHalf, TryGet[uint8](retriever))
	}))

	for i := 0; i < 255; i++ {
		expected := fmt.Sprintf("%s %d", strFirstHalf, i+1)
		if result := TryGet[string](subInjector); result != expected {
			t.Errorf("sub-injector is not calling parent injector correctly. Received: `%s`; Expected: `%s`", result, expected)
			t.FailNow()
		}
	}
}

func TestStdInjector_SubInjector__OverrideParent(t *testing.T) {
	const strFirstHalf = "The totally value of it is"
	parent := New(false, types.ReflectionOptions{})
	subInjector := parent.SubInjector(false)

	Register(parent, binds.Factory(func(retriever types.DependencyRetriever) uint8 {
		return 101
	}))

	Register(subInjector, binds.Factory(func(retriever types.DependencyRetriever) string {
		return fmt.Sprintf("%s %d", strFirstHalf, TryGet[uint8](retriever))
	}))

	expected := fmt.Sprintf("%s 101", strFirstHalf)
	if result := TryGet[string](subInjector); result != expected {
		t.Errorf("sub-injector is not calling parent injector correctly. Received: `%s`; Expected: `%s`", result, expected)
		t.FailNow()
	}

	// Register a new uint8 to override parent
	Register(subInjector, binds.Singleton(func(retriever types.DependencyRetriever) uint8 {
		return 42
	}))

	expected = fmt.Sprintf("%s 42", strFirstHalf)
	if result := TryGet[string](subInjector); result != expected {
		t.Errorf("sub-injector is not calling parent injector correctly. Received: `%s`; Expected: `%s`", result, expected)
		t.FailNow()
	}

	// Checks if parent still returns the same old value
	parentResult := TryGet[uint8](parent)
	if parentResult != 101 {
		t.Errorf("parent value was overrided, it should not occur")
	}
}
