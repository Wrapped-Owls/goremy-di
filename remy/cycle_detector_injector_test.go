package remy

import (
	"errors"
	"fmt"
	"testing"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
)

func TestCycleDetectorInjector_Register(t *testing.T) {
	defer func() {
		r := recover()

		asErr, ok := r.(error)
		if !ok {
			t.Fatalf("Register() did not return an error")
		}
		if r != nil && !errors.Is(asErr, remyErrs.ErrCycleDependencyDetectedSentinel) {
			t.Error(r)
		}
	}()
	ij := NewCycleDetectorInjector(Config{CanOverride: false})
	cycleKey := [...]string{"lang", "tool"}
	Register(
		ij, Factory(
			func(retriever DependencyRetriever) (result string, err error) {
				result = MustGet[string](retriever, cycleKey[0]) + " is awesome"
				return
			},
		),
	)
	Register(
		ij, Factory(
			func(retriever DependencyRetriever) (result string, err error) {
				result = "git" + MustGet[string](retriever)
				return
			},
		), cycleKey[1],
	)
	Register(
		ij, Factory(
			func(retriever DependencyRetriever) (result string, err error) {
				result = "Go + " + MustGet[string](retriever, cycleKey[1])
				return
			},
		), cycleKey[0],
	)
	Register(
		ij, Singleton(
			func(retriever DependencyRetriever) (result int, err error) {
				result = len(MustGet[string](retriever, cycleKey[0]))
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
					MustGet[string](retriever), MustGet[uint8](retriever),
				)
				return
			},
		), cycleKey,
	)

	if _, err := Get[string](ij, cycleKey); err != nil {
		t.Errorf("Something went wrong during normal utilization, raise: %v", err)
	}

	// overrides a dependency to insert a cycle
	Override(
		ij, Factory(
			func(retriever DependencyRetriever) (uint8, error) {
				val, err := Get[string](retriever, cycleKey)
				if err != nil {
					return 0, err
				}
				return uint8(len(val)), nil
			},
		),
	)
	_, err := Get[string](ij, cycleKey)
	if err == nil {
		t.Error("function executes normally when it should raise an error")
		t.FailNow()
	}

	if !errors.Is(err, ErrCycleDependencyDetected) {
		t.Errorf("The returned error is not ErrCycleDependencyDetected")
	}
}
