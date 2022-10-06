package remy

import (
	"fmt"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

func TestCycleDetectorInjector_Register(t *testing.T) {
	defer func() {
		r := recover()
		if r != nil && fmt.Sprint(r) != utils.ErrCycleDependencyDetected.Error() {
			t.Error(r)
		}
	}()
	ij := NewCycleDetectorInjector(Config{CanOverride: false})
	cycleKey := [...]string{"lang", "tool"}
	Register(
		ij, Factory(
			func(retriever DependencyRetriever) (result string, err error) {
				result = Get[string](retriever, cycleKey[0]) + " is awesome"
				return
			},
		),
	)
	Register(
		ij, Factory(
			func(retriever DependencyRetriever) (result string, err error) {
				result = "git" + Get[string](retriever)
				return
			},
		), cycleKey[1],
	)
	Register(
		ij, Factory(
			func(retriever DependencyRetriever) (result string, err error) {
				result = "Go + " + Get[string](retriever, cycleKey[1])
				return
			},
		), cycleKey[0],
	)
	Register(
		ij, Singleton(
			func(retriever DependencyRetriever) (result int, err error) {
				result = len(Get[string](retriever, cycleKey[0]))
				return
			},
		),
	)
}

func TestCycleDetectorInjector_Get(t *testing.T) {
	ij := NewCycleDetectorInjector(Config{CanOverride: true})
	const cycleKey = "name"
	RegisterInstance(ij, "go")
	RegisterInstance(ij, uint8(42))
	Register(
		ij, Factory(
			func(retriever DependencyRetriever) (result string, err error) {
				result = fmt.Sprintf(
					"The lenght for the string `%s` is %d ",
					Get[string](retriever), Get[uint8](retriever),
				)
				return
			},
		), cycleKey,
	)

	if _, err := DoGet[string](ij, cycleKey); err != nil {
		t.Errorf("Something went wrong during normal utilization, raise: %v", err)
	}

	// overrides a dependency to insert a cycle
	Override(
		ij, Factory(
			func(retriever DependencyRetriever) (uint8, error) {
				return uint8(len(Get[string](retriever, cycleKey))), nil
			},
		),
	)
	_, err := DoGet[string](ij, cycleKey)
	if err == nil {
		t.Error("function executes normally when it should raise an error")
	}
}
