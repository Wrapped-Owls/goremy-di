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
