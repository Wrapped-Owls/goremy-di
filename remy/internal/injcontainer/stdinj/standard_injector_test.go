package stdinj

import (
	"errors"
	"reflect"
	"testing"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
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

func TestStdInjector_BindElemOverride(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         string
		cache        injopts.CacheConfOption
		softOverride bool
		tag          string
		wantErr      error
	}{
		{
			name:    "rebinding is refused by default",
			cache:   injopts.CacheOptNone,
			wantErr: remyErrs.ErrAlreadyBoundSentinel,
		},
		{
			name:         "an allowed scope still refuses an unflagged rebind",
			cache:        injopts.CacheOptAllowOverride,
			softOverride: false,
			wantErr:      remyErrs.ErrAlreadyBoundSentinel,
		},
		{
			name:         "a flagged rebind on an allowed scope succeeds",
			cache:        injopts.CacheOptAllowOverride,
			softOverride: true,
			wantErr:      nil,
		},
		{
			name:         "the flag alone is not enough",
			cache:        injopts.CacheOptNone,
			softOverride: true,
			wantErr:      remyErrs.ErrAlreadyBoundSentinel,
		},
		{
			name:         "a tagged key follows the same rule",
			cache:        injopts.CacheOptNone,
			softOverride: true,
			tag:          "edition",
			wantErr:      remyErrs.ErrAlreadyBoundSentinel,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			inj := New(Options{Cache: testCase.cache})
			key := types.KeyElem[string]{}
			first := types.BindOptions{Tag: testCase.tag}
			if err := inj.BindElem(key, "first", first); err != nil {
				t.Fatalf("first BindElem: %v", err)
			}

			rebind := types.BindOptions{Tag: testCase.tag, SoftOverride: testCase.softOverride}
			err := inj.BindElem(key, "second", rebind)
			if !errors.Is(err, testCase.wantErr) {
				t.Fatalf("second BindElem = %v, want %v", err, testCase.wantErr)
			}
		})
	}
}

func TestStdInjector_GetAll(t *testing.T) {
	t.Parallel()

	const listing = injopts.CacheOptReturnAll

	testCases := []struct {
		name             string
		parentOpts       Options
		childOpts        Options
		tag              string
		want             []any
		wantErr          error
		wantParentBlamed bool
	}{
		{
			name:       "a scope lists its own elements before its parent's",
			parentOpts: Options{Cache: listing},
			childOpts:  Options{Cache: listing},
			want:       []any{42, "from-parent"},
		},
		{
			name:       "a tagged listing merges the same way",
			parentOpts: Options{Cache: listing},
			childOpts:  Options{Cache: listing},
			tag:        "edition",
			want:       []any{42, "from-parent"},
		},
		{
			name:       "a scope that may not list still gets its parent's",
			parentOpts: Options{Cache: listing},
			childOpts:  Options{},
			want:       []any{"from-parent"},
		},
		{
			name:       "an isolated scope never reads the parent",
			parentOpts: Options{Cache: listing},
			childOpts:  Options{Cache: listing, Resolve: injopts.ResolveOptIsolated},
			want:       []any{42},
		},
		{
			name:             "a parent that may not list fails the whole listing",
			parentOpts:       Options{},
			childOpts:        Options{Cache: listing},
			wantErr:          remyErrs.ErrConfigNotAllowReturnAll,
			wantParentBlamed: true,
		},
		{
			name:       "when neither may list, the scope reports its own refusal",
			parentOpts: Options{},
			childOpts:  Options{},
			wantErr:    remyErrs.ErrConfigNotAllowReturnAll,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			parent := New(testCase.parentOpts)
			child := New(testCase.childOpts, parent)
			store(t, parent, testCase.tag, "from-parent")
			store(t, child, testCase.tag, 42)

			got, err := child.GetAll(testCase.tag)
			if !errors.Is(err, testCase.wantErr) {
				t.Fatalf("GetAll error = %v, want %v", err, testCase.wantErr)
			}

			var fromParent remyErrs.ErrWrapParentSubErrors
			if blamesParent := errors.As(
				err,
				&fromParent,
			); blamesParent != testCase.wantParentBlamed {
				t.Fatalf(
					"blamed the parent = %v, want %v (err %v)",
					blamesParent,
					testCase.wantParentBlamed,
					err,
				)
			}
			if testCase.wantErr != nil {
				return
			}

			if !reflect.DeepEqual(got, testCase.want) {
				t.Fatalf("GetAll = %#v, want %#v", got, testCase.want)
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
