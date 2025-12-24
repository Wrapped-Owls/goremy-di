package injector

import (
	"fmt"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

func TestStdInjector_SubInjector(t *testing.T) {
	const strFirstHalf = "the counter is at"
	parent := New(injopts.CacheOptNone)
	subInjector := parent.SubInjector(false)

	var counter uint8 = 0
	_ = Register(
		parent, "", binds.Factory(
			func(retriever types.DependencyRetriever) (uint8, error) {
				counter++
				return counter, nil
			},
		),
	)

	_ = Register(
		subInjector, "", binds.Factory(
			func(retriever types.DependencyRetriever) (string, error) {
				return fmt.Sprintf("%s %d", strFirstHalf, TryGet[uint8](retriever, "")), nil
			},
		),
	)

	for i := 0; i < 255; i++ {
		expected := fmt.Sprintf("%s %d", strFirstHalf, i+1)
		if result := TryGet[string](subInjector, ""); result != expected {
			t.Errorf(
				"sub-injector is not calling parent injector correctly. Received: `%s`; Expected: `%s`",
				result,
				expected,
			)
			t.FailNow()
		}
	}
}

func TestStdInjector_SubInjectorEmpty(t *testing.T) {
	const elementKey = "game-name"
	parent := New(injopts.CacheOptNone)
	subInjector := parent.SubInjector(false)

	_ = Register(parent, elementKey, binds.Instance("snake-pong"))

	results := [...]string{
		TryGet[string](parent, elementKey),
		TryGet[string](subInjector, elementKey),
	}
	if results[0] != results[1] {
		t.Error("Result isn't the same for parent and sub injectors")
	}
}

func TestStdInjector_GetUnboundedElement(t *testing.T) {
	const errMessage = "An error have not been returned when getting unbounded element"
	parentInjector := New(injopts.CacheOptNone)
	for _, ij := range [...]types.Injector{parentInjector, parentInjector.SubInjector()} {
		if _, err := Get[string](ij, ""); err == nil {
			t.Error(errMessage)
		}
		if _, err := Get[uint8](ij, "release-date"); err == nil {
			t.Error(errMessage)
		}
	}
}

func TestStdInjector_SubInjector__OverrideParent(t *testing.T) {
	const strFirstHalf = "The totally value of it is"
	parent := New(injopts.CacheOptNone)
	subInjector := parent.SubInjector(false)

	_ = Register(
		parent, "", binds.Factory(
			func(retriever types.DependencyRetriever) (uint8, error) {
				return 101, nil
			},
		),
	)

	_ = Register(
		subInjector, "", binds.Factory(
			func(retriever types.DependencyRetriever) (string, error) {
				return fmt.Sprintf("%s %d", strFirstHalf, TryGet[uint8](retriever, "")), nil
			},
		),
	)

	expected := fmt.Sprintf("%s 101", strFirstHalf)
	if result := TryGet[string](subInjector, ""); result != expected {
		t.Errorf(
			"sub-injector is not calling parent injector correctly. Received: `%s`; Expected: `%s`",
			result, expected,
		)
		t.FailNow()
	}

	// Register a new uint8 to override parent
	_ = Register(
		subInjector, "", binds.Singleton(
			func(retriever types.DependencyRetriever) (uint8, error) {
				return 42, nil
			},
		),
	)

	expected = fmt.Sprintf("%s 42", strFirstHalf)
	if result := TryGet[string](subInjector, ""); result != expected {
		t.Errorf(
			"sub-injector is not calling parent injector correctly. Received: `%s`; Expected: `%s`",
			result, expected,
		)
		t.FailNow()
	}

	// Checks if parent still returns the same old value
	parentResult := TryGet[uint8](parent, "")
	if parentResult != 101 {
		t.Errorf("parent value was overrided, it should not occur")
	}
}
