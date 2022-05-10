package injector

import (
	"fmt"
	"github.com/wrapped-owls/fitpiece/gotalaria/internal/binds"
	"github.com/wrapped-owls/fitpiece/gotalaria/internal/types"
	"testing"
)

func TestStdInjector_SubInjector(t *testing.T) {
	const strFirstHalf = "the counter is at"
	parent := New(false)
	subInjector := parent.SubInjector(false)

	var counter uint8 = 0
	Register[uint8](parent, binds.Factory(func(retriever types.DependencyRetriever) uint8 {
		counter++
		return counter
	}))

	Register[string](subInjector, binds.Factory(func(retriever types.DependencyRetriever) string {
		return fmt.Sprintf("%s %d", strFirstHalf, Get[uint8](retriever))
	}))

	for i := 0; i < 255; i++ {
		expected := fmt.Sprintf("%s %d", strFirstHalf, i+1)
		if result := Get[string](subInjector); result != expected {
			t.Errorf("sub-injector is not calling parent injector correctly. Received: `%s`; Expected: `%s`", result, expected)
			t.FailNow()
		}
	}
}

func TestStdInjector_SubInjector__OverrideParent(t *testing.T) {
	const strFirstHalf = "The totally value of it is"
	parent := New(false)
	subInjector := parent.SubInjector(false)

	Register[uint8](parent, binds.Factory(func(retriever types.DependencyRetriever) uint8 {
		return 101
	}))

	Register[string](subInjector, binds.Factory(func(retriever types.DependencyRetriever) string {
		return fmt.Sprintf("%s %d", strFirstHalf, Get[uint8](retriever))
	}))

	expected := fmt.Sprintf("%s 101", strFirstHalf)
	if result := Get[string](subInjector); result != expected {
		t.Errorf("sub-injector is not calling parent injector correctly. Received: `%s`; Expected: `%s`", result, expected)
		t.FailNow()
	}

	// Register a new uint8 to override parent
	Register[uint8](subInjector, binds.Singleton(func(retriever types.DependencyRetriever) uint8 {
		return 42
	}))

	expected = fmt.Sprintf("%s 42", strFirstHalf)
	if result := Get[string](subInjector); result != expected {
		t.Errorf("sub-injector is not calling parent injector correctly. Received: `%s`; Expected: `%s`", result, expected)
		t.FailNow()
	}

	// Checks if parent still returns the same old value
	parentResult := Get[uint8](parent)
	if parentResult != 101 {
		t.Errorf("parent value was overrided, it should not occur")
	}
}
