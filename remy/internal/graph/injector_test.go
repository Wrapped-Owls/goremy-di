package graph

import (
	"errors"
	"runtime"
	"strings"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/injcontainer"
	"github.com/wrapped-owls/goremy-di/remy/internal/injcontainer/stdinj"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type leafDep struct{ value string }

// numberedKey gives each level of a deep chain a distinct BindKey
func numberedKey(level int) types.BindKey {
	keys := []types.BindKey{
		types.KeyElem[int8]{},
		types.KeyElem[int16]{},
		types.KeyElem[int32]{},
		types.KeyElem[int64]{},
		types.KeyElem[uint8]{},
		types.KeyElem[uint16]{},
		types.KeyElem[uint32]{},
		types.KeyElem[float32]{},
		types.KeyElem[float64]{},
		types.KeyElem[complex64]{},
		types.KeyElem[complex128]{},
		types.KeyElem[uintptr]{},
	}
	return keys[level%len(keys)]
}

func newBase() types.Injector {
	return stdinj.New(stdinj.Options{})
}

// registerChain wires string -> int -> bool, each resolved lazily
func registerChain(t *testing.T, inj types.Injector) {
	t.Helper()

	err := errors.Join(
		injcontainer.Register[string](inj, "", binds.Factory(
			func(r types.DependencyRetriever) (string, error) {
				value, getErr := injcontainer.Get[int](r, "")
				return string(rune(value)), getErr
			},
		)),
		injcontainer.Register[int](inj, "", binds.Factory(
			func(r types.DependencyRetriever) (int, error) {
				_, getErr := injcontainer.Get[bool](r, "")
				return 65, getErr
			},
		)),
		injcontainer.Register[bool](inj, "", binds.Instance(true)),
	)
	if err != nil {
		t.Fatalf("register chain: %v", err)
	}
}

func TestInjector_RecordsEdges(t *testing.T) {
	t.Parallel()

	inj, graph := New(newBase())
	registerChain(t, inj)

	if _, err := injcontainer.Get[string](inj, ""); err != nil {
		t.Fatalf("unexpected resolution error: %v", err)
	}

	edges := graph.Edges()
	if len(edges) != 2 {
		t.Fatalf("edges = %d, want 2 (string->int, int->bool): %v", len(edges), edges)
	}

	wantPairs := map[[2]types.BindKey]bool{
		{types.KeyElem[string]{}, types.KeyElem[int]{}}: false,
		{types.KeyElem[int]{}, types.KeyElem[bool]{}}:   false,
	}
	for _, edge := range edges {
		pair := [2]types.BindKey{edge.From.Key, edge.To.Key}
		if _, expected := wantPairs[pair]; !expected {
			t.Errorf("unexpected edge: %v", edge)
			continue
		}
		wantPairs[pair] = true
	}
	for pair, found := range wantPairs {
		if !found {
			t.Errorf("missing edge %v", pair)
		}
	}
}

func TestInjector_CycleRendersOrderedPath(t *testing.T) {
	t.Parallel()

	inj := NewDetector(newBase())

	// string -> int -> bool -> string closes a cycle
	registerErr := errors.Join(
		injcontainer.Register[string](inj, "", binds.Factory(
			func(r types.DependencyRetriever) (string, error) {
				_, err := injcontainer.Get[int](r, "")
				return "", err
			},
		)),
		injcontainer.Register[int](inj, "", binds.Factory(
			func(r types.DependencyRetriever) (int, error) {
				_, err := injcontainer.Get[bool](r, "")
				return 0, err
			},
		)),
		injcontainer.Register[bool](inj, "", binds.Factory(
			func(r types.DependencyRetriever) (bool, error) {
				_, err := injcontainer.Get[string](r, "")
				return false, err
			},
		)),
	)
	if registerErr != nil {
		t.Fatalf("register cycle: %v", registerErr)
	}

	errMsg := resolvePanicMessage(t, func() { _, _ = injcontainer.Get[string](inj, "") })
	positions := []int{
		strings.Index(errMsg, "string"),
		strings.Index(errMsg, "int"),
		strings.Index(errMsg, "bool"),
		strings.LastIndex(errMsg, "string"),
	}
	for index := 1; index < len(positions); index++ {
		if positions[index-1] < 0 || positions[index-1] >= positions[index] {
			t.Fatalf("cycle path is not ordered string -> int -> bool -> string: %q", errMsg)
		}
	}
}

func TestInjector_DiamondIsNotCycle(t *testing.T) {
	t.Parallel()

	inj, graph := New(newBase())

	// string and int both depend on leafDep, and bool depends on both
	registerErr := errors.Join(
		injcontainer.Register[leafDep](inj, "", binds.Instance(leafDep{value: "shared"})),
		injcontainer.Register[string](inj, "", binds.Factory(
			func(r types.DependencyRetriever) (string, error) {
				leaf, err := injcontainer.Get[leafDep](r, "")
				return leaf.value, err
			},
		)),
		injcontainer.Register[int](inj, "", binds.Factory(
			func(r types.DependencyRetriever) (int, error) {
				leaf, err := injcontainer.Get[leafDep](r, "")
				return len(leaf.value), err
			},
		)),
		injcontainer.Register[bool](inj, "", binds.Factory(
			func(r types.DependencyRetriever) (bool, error) {
				if _, err := injcontainer.Get[string](r, ""); err != nil {
					return false, err
				}
				_, err := injcontainer.Get[int](r, "")
				return true, err
			},
		)),
	)
	if registerErr != nil {
		t.Fatalf("register diamond: %v", registerErr)
	}

	if _, err := injcontainer.Get[bool](inj, ""); err != nil {
		t.Fatalf("diamond resolution must not be a cycle: %v", err)
	}
	if edges := graph.Edges(); len(edges) != 4 {
		t.Fatalf("edges = %d, want 4: %v", len(edges), edges)
	}
}

// resolvePanicMessage runs resolve and returns the cycle error it panics with
func resolvePanicMessage(t *testing.T, resolve func()) string {
	t.Helper()

	var message string
	func() {
		defer func() {
			recovered := recover()
			asErr, ok := recovered.(error)
			if !ok {
				t.Fatalf("expected a cycle panic, got %v", recovered)
			}
			if !errors.Is(asErr, remyErrs.ErrCycleDependencyDetectedSentinel) {
				t.Fatalf("expected cycle error, got: %v", asErr)
			}
			message = asErr.Error()
		}()
		resolve()
	}()
	return message
}

func TestInjector_CycleBeyondInlineCapacity(t *testing.T) {
	t.Parallel()

	head := types.BindKey(types.KeyElem[string]{})
	deep := types.Injector(NewDetector(newBase())).RetrieverFor(head, "")
	for level := 0; level < inlinePathCapacity+2; level++ {
		deep = deep.RetrieverFor(numberedKey(level), "")
	}

	defer func() {
		recovered := recover()
		asErr, ok := recovered.(error)
		if !ok {
			t.Fatalf("revisiting the head must panic, got %v", recovered)
		}
		if !errors.Is(asErr, remyErrs.ErrCycleDependencyDetectedSentinel) {
			t.Fatalf("expected a cycle error, got %v", asErr)
		}
		if path := asErr.(*remyErrs.ErrCycleDependencyDetected).Path; len(
			path,
		) != inlinePathCapacity+4 {
			t.Fatalf("cycle path has %d nodes, want the whole deep chain", len(path))
		}
	}()
	deep.RetrieverFor(head, "")
}

func TestInjector_SiblingsBeyondInlineCapacityStayIndependent(t *testing.T) {
	t.Parallel()

	root := types.Injector(NewDetector(newBase())).RetrieverFor(types.KeyElem[string]{}, "")

	deep := root
	for level := 0; level < inlinePathCapacity+2; level++ {
		deep = deep.RetrieverFor(numberedKey(level), "")
	}

	// numberedKey(0) is on the deep chain but not on the sibling one, so forking
	// it from the root again must be a diamond, never a cycle
	if sibling := root.RetrieverFor(numberedKey(0), ""); sibling == nil {
		t.Fatal("sibling fork returned nil")
	}
}

func BenchmarkInjector_RetrieverFor(b *testing.B) {
	inj := NewDetector(newBase())
	if err := injcontainer.Register[int](inj, "", binds.Instance(42)); err != nil {
		b.Fatalf("register: %v", err)
	}

	var sink int
	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		sink, _ = injcontainer.Get[int](inj, "")
	}
	runtime.KeepAlive(sink)
}

func BenchmarkInjector_SubInjector(b *testing.B) {
	inj := NewDetector(newBase())

	var sink types.Injector
	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		sink = inj.SubInjector()
	}
	runtime.KeepAlive(sink)
}
