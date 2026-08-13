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
	parent := New(Options{})
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
	parent := New(Options{})
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
	parentInjector := New(Options{})
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
	parent := New(Options{})
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

func TestStdInjector_ResolveOptions(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		parent *StdInjector
		child  Options
		want   injopts.ResolveConfOption
	}{
		{
			name:  "own options are reported",
			child: Options{Resolve: injopts.ResolveOptDuckTyping},
			want:  injopts.ResolveOptDuckTyping,
		},
		{
			name:   "inheritable options come down from the parent",
			parent: New(Options{Resolve: injopts.ResolveOptTracePath}),
			child:  Options{},
			want:   injopts.ResolveOptTracePath,
		},
		{
			name:   "own and inherited options merge",
			parent: New(Options{Resolve: injopts.ResolveOptTracePath}),
			child:  Options{Resolve: injopts.ResolveOptDuckTyping},
			want:   injopts.ResolveOptDuckTyping | injopts.ResolveOptTracePath,
		},
		{
			name:   "an isolated scope inherits nothing",
			parent: New(Options{Resolve: injopts.ResolveOptTracePath}),
			child:  Options{Resolve: injopts.ResolveOptIsolated},
			want:   injopts.ResolveOptIsolated,
		},
		{
			name:   "isolation itself never propagates to a child",
			parent: New(Options{Resolve: injopts.ResolveOptIsolated}),
			child:  Options{},
			want:   injopts.ResolveOptNone,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			var inj *StdInjector
			if testCase.parent != nil {
				inj = New(testCase.child, testCase.parent)
			} else {
				inj = New(testCase.child)
			}

			if got := inj.ResolveOptions(); got != testCase.want {
				t.Fatalf("ResolveOptions() = %b, want %b", got, testCase.want)
			}
		})
	}
}

func TestStdInjector_RetrieverForOptsOut(t *testing.T) {
	t.Parallel()

	// the default injector does not track resolutions, so it keeps the caller's
	// retriever instead of allocating a scoped one
	if scoped := New(Options{}).RetrieverFor(types.KeyElem[string]{}, ""); scoped != nil {
		t.Fatalf("RetrieverFor = %v, want nil on the non-tracking injector", scoped)
	}
}
