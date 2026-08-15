package stdinj

import (
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

// the container stores whatever it is handed and gives it back untouched, so plain
// values keep these assertions about the container instead of about binds
func store[T any](tb testing.TB, inj types.Injector, tag string, value T) {
	tb.Helper()

	if err := inj.BindElem(types.KeyElem[T]{}, value, types.BindOptions{Tag: tag}); err != nil {
		tb.Fatalf("BindElem: %v", err)
	}
}

func fetch[T any](tb testing.TB, retriever types.DependencyRetriever, tag string) (result T) {
	tb.Helper()

	stored, err := retriever.RetrieveBind(types.KeyElem[T]{}, tag)
	if err != nil {
		tb.Fatalf("RetrieveBind: %v", err)
	}

	value, ok := stored.(T)
	if !ok {
		tb.Fatalf("stored value is %T, want %T", stored, result)
	}
	return value
}

func TestStdInjector_SubInjectorFallsBackToParent(t *testing.T) {
	t.Parallel()

	parent := New(Options{})
	subInjector := parent.SubInjector(false)

	var calls int
	store(t, parent, "", func() int { calls++; return calls })
	store(t, parent, "game-name", "snake-pong")

	// a function proves identity: the container hands back the very closure it was
	// given, so calling it through the sub-injector advances the parent's counter
	for round := 1; round <= 3; round++ {
		if got := fetch[func() int](t, subInjector, "")(); got != round {
			t.Fatalf("call %d through the sub-injector returned %d, want %d", round, got, round)
		}
	}

	if got := fetch[string](t, subInjector, "game-name"); got != "snake-pong" {
		t.Fatalf("tagged lookup through the sub-injector = %q, want %q", got, "snake-pong")
	}
}

func TestStdInjector_SubInjectorShadowsParent(t *testing.T) {
	t.Parallel()

	parent := New(Options{})
	subInjector := parent.SubInjector(false)
	store(t, parent, "", 101)

	if got := fetch[int](t, subInjector, ""); got != 101 {
		t.Fatalf("before shadowing = %d, want the parent value 101", got)
	}

	store(t, subInjector, "", 42)

	if got := fetch[int](t, subInjector, ""); got != 42 {
		t.Fatalf("after shadowing = %d, want the sub-injector value 42", got)
	}
	if got := fetch[int](t, parent, ""); got != 101 {
		t.Fatalf("parent = %d, want 101: shadowing must not reach upwards", got)
	}
}

func TestStdInjector_RetrieveBindUnbound(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		tag  string
	}{
		{name: "untagged key was never bound", tag: ""},
		{name: "tagged key was never bound", tag: "release-date"},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			parent := New(Options{})
			for _, inj := range [...]types.Injector{parent, parent.SubInjector()} {
				if _, err := inj.RetrieveBind(types.KeyElem[string]{}, testCase.tag); err == nil {
					t.Errorf("RetrieveBind on %T returned no error for an unbound key", inj)
				}
			}
		})
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
