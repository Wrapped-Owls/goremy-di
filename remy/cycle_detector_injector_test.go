package remy

import (
	"fmt"
	"testing"
)

func TestCycleDetectorInjector_Get(t *testing.T) {
	ij := NewCycleDetectorInjector(Config{CanOverride: true})
	const cycleKey = "name"
	RegisterInstance(ij, "go")
	RegisterInstance(ij, uint8(42))
	Register(ij, Factory(func(retriever DependencyRetriever) string {
		return fmt.Sprintf(
			"The lenght for the string `%s` is %d ",
			Get[string](retriever), Get[uint8](retriever),
		)
	}), cycleKey)

	if _, err := DoGet[string](ij, cycleKey); err != nil {
		t.Errorf("Something went wrong during normal utilization, raise: %v", err)
	}

	// overrides a dependency to insert a cycle
	Override(ij, Factory(func(retriever DependencyRetriever) uint8 {
		return uint8(len(Get[string](retriever, cycleKey)))
	}))
	_, err := DoGet[string](ij, cycleKey)
	if err == nil {
		t.Error("function executes normally when it should raise an error")
	}
}
