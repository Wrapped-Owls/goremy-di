package remy

import (
	"errors"
	"fmt"
	"testing"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
)

func TestNewGraphInjector(t *testing.T) {
	inj, graph := NewGraphInjector()
	Register(inj, Factory(func(r DependencyRetriever) (string, error) {
		value, err := Get[int](r)
		return string(rune(value)), err
	}))
	RegisterInstance(inj, 65)

	if _, err := Get[string](inj); err != nil {
		t.Fatalf("unexpected resolution error: %v", err)
	}

	edges := graph.Edges()
	if len(edges) != 1 {
		t.Fatalf("edges = %d, want the string -> int pair: %v", len(edges), edges)
	}
	if edges[0].From.Key != NewBindKey[string]() || edges[0].To.Key != NewBindKey[int]() {
		t.Fatalf("unexpected edge: %v", edges[0])
	}
}

func TestNewCycleDetector(t *testing.T) {
	inj := NewCycleDetector()
	Register(inj, Factory(func(r DependencyRetriever) (string, error) {
		_, err := Get[int](r)
		return "", err
	}))
	Register(inj, Factory(func(r DependencyRetriever) (int, error) {
		_, err := Get[string](r)
		return 0, err
	}))

	if _, err := Get[string](inj); !errors.Is(err, ErrCycleDependencyDetected) {
		t.Fatalf("expected cycle error, got: %v", err)
	}
}

func TestNewGraphInjector_HonoursParentInjector(t *testing.T) {
	parent := NewInjector()
	RegisterInstance(parent, "from-parent")

	inj, _ := NewGraphInjector(Config{ParentInjector: parent})
	if got := MustGet[string](inj); got != "from-parent" {
		t.Fatalf("Get = %q, want the parent bind", got)
	}
}

func TestNewCycleDetector_HonoursParentInjector(t *testing.T) {
	parent := NewInjector()
	RegisterInstance(parent, "from-parent")

	inj := NewCycleDetector(Config{ParentInjector: parent})
	if got := MustGet[string](inj); got != "from-parent" {
		t.Fatalf("Get = %q, want the parent bind", got)
	}
}

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
	ij := NewCycleDetector(Config{CanOverride: false})
	cycleKey := [...]Tag{"lang", "tool"}
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

func TestCycleDetectorInjector_RegisterTimeSelfCycle(t *testing.T) {
	defer func() {
		recovered := recover()
		asErr, ok := recovered.(error)
		if !ok {
			t.Fatalf("Register() should panic with an error, got %v", recovered)
		}
		if !errors.Is(asErr, remyErrs.ErrCycleDependencyDetectedSentinel) {
			t.Errorf("expected cycle detection, got: %v", asErr)
		}
	}()

	ij := NewCycleDetector()
	Register(
		ij, Factory(
			func(retriever DependencyRetriever) (int, error) {
				value, err := Get[uint8](retriever)
				return int(value), err
			},
		),
	)
	// registering this singleton generates it immediately; its chain circles
	// back into uint8 (the key being registered), which must be a cycle error,
	// not an element-not-registered one
	Register(
		ij, Singleton(
			func(retriever DependencyRetriever) (uint8, error) {
				value, err := Get[int](retriever)
				return uint8(value), err
			},
		),
	)
}

func TestCycleDetectorInjector_Get(t *testing.T) {
	ij := NewCycleDetector(Config{CanOverride: true})
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

//nolint:staticcheck // exercises the deprecated wrapper on purpose
func TestNewCycleDetectorInjector_DelegatesToDetector(t *testing.T) {
	inj := NewCycleDetectorInjector()
	Register(inj, Factory(func(r DependencyRetriever) (string, error) {
		_, err := Get[int](r)
		return "", err
	}))
	Register(inj, Factory(func(r DependencyRetriever) (int, error) {
		_, err := Get[string](r)
		return 0, err
	}))

	if _, err := Get[string](inj); !errors.Is(err, ErrCycleDependencyDetected) {
		t.Fatalf("deprecated wrapper lost cycle detection: %v", err)
	}
}

type (
	cycleHead struct{ depth int }
	cycleTail struct{ depth int }
)

func TestNewCycleDetector_CycleCrossingATemporaryScope(t *testing.T) {
	testCases := []struct {
		name  string
		reach func(retriever DependencyRetriever) (cycleTail, error)
	}{
		{
			name: "through GetWith",
			reach: func(retriever DependencyRetriever) (cycleTail, error) {
				return GetWith[cycleTail](retriever, func(Injector) error { return nil })
			},
		},
		{
			name: "through GetWithPairs",
			reach: func(retriever DependencyRetriever) (cycleTail, error) {
				return GetWithPairs[cycleTail](retriever, []BindEntry{NewBindEntry("scoped")})
			},
		},
		{
			name: "through a sub-injector taken mid-resolution",
			reach: func(retriever DependencyRetriever) (cycleTail, error) {
				injector, isInjector := retriever.(Injector)
				if !isInjector {
					return cycleTail{}, errors.New("retriever is not an Injector")
				}
				return Get[cycleTail](injector.SubInjector())
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			inj := NewCycleDetector()
			RegisterFactory(inj, func(retriever DependencyRetriever) (cycleHead, error) {
				tail, err := testCase.reach(retriever)
				return cycleHead{depth: tail.depth + 1}, err
			})
			RegisterFactory(inj, func(retriever DependencyRetriever) (cycleTail, error) {
				head, err := Get[cycleHead](retriever)
				return cycleTail{depth: head.depth + 1}, err
			})

			_, err := Get[cycleHead](inj)
			if !errors.Is(err, ErrCycleDependencyDetected) {
				t.Fatalf("err = %v, want a reported cycle", err)
			}
		})
	}
}
