package binds

import (
	"errors"
	"runtime"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

func suffixDecorator(_ types.DependencyRetriever, value string) (string, error) {
	return value + "-decorated", nil
}

func TestDecorate_Generates(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		bind types.Bind[string]
		want types.BindType
	}{
		{name: "Instance", bind: Instance("plain"), want: types.BindInstance},
		{name: "Factory", bind: Factory(constantBinder("plain")), want: types.BindFactory},
		{name: "Singleton", bind: Singleton(constantBinder("plain")), want: types.BindSingleton},
		{
			name: "LazySingleton",
			bind: LazySingleton(constantBinder("plain")),
			want: types.BindLazySingleton,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			decorated := Decorate[string](testCase.bind, suffixDecorator)

			got, err := decorated.Generates(nil)
			if err != nil {
				t.Fatalf("Generates() error = %v, want nil", err)
			}
			if got != "plain-decorated" {
				t.Fatalf("Generates() = %q, want the wrapped value decorated", got)
			}
			if kind := decorated.Type(); kind != testCase.want {
				t.Fatalf("Type() = %d, want the wrapped kind %d", kind, testCase.want)
			}
		})
	}
}

func TestDecorate_FoldsTheDecorationIntoTheBindItself(t *testing.T) {
	t.Parallel()

	for _, bind := range []types.Bind[string]{
		Instance("plain"),
		Factory(constantBinder("plain")),
		Singleton(constantBinder("plain")),
		LazySingleton(constantBinder("plain")),
	} {
		if _, wrapped := Decorate[string](bind, suffixDecorator).(DecoratedBind[string]); wrapped {
			t.Fatalf("%T was wrapped instead of rebuilt around the decorated generator", bind)
		}
	}
}

func TestDecorate_KeepsTheWrappedBindCaching(t *testing.T) {
	t.Parallel()

	var calls int
	decorated := Decorate[int](LazySingleton(countingBinder(&calls)), func(
		_ types.DependencyRetriever, value int,
	) (int, error) {
		return value * 10, nil
	})

	for round := 1; round <= 3; round++ {
		got, err := decorated.Generates(nil)
		if err != nil {
			t.Fatalf("call %d: unexpected error: %v", round, err)
		}
		if got != 10 {
			t.Fatalf("call %d = %d, want the one built value decorated", round, got)
		}
	}
	if calls != 1 {
		t.Fatalf("binder ran %d times, want the singleton to build once", calls)
	}
}

func TestDecorate_PropagatesErrors(t *testing.T) {
	t.Parallel()

	errDecorator := errors.New("decorator failed")

	testCases := []struct {
		name    string
		bind    types.Bind[string]
		wantErr error
	}{
		{
			name:    "the wrapped bind fails before the decorator runs",
			bind:    Decorate[string](Factory(failingBinder[string]()), suffixDecorator),
			wantErr: errBinderFailed,
		},
		{
			name: "the decorator itself fails",
			bind: Decorate[string](Factory(constantBinder("plain")), func(
				types.DependencyRetriever, string,
			) (string, error) {
				return "", errDecorator
			}),
			wantErr: errDecorator,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if _, err := testCase.bind.Generates(nil); !errors.Is(err, testCase.wantErr) {
				t.Fatalf("Generates() error = %v, want %v", err, testCase.wantErr)
			}
		})
	}
}

func TestDecorate_KeepsTheDuckTypingViews(t *testing.T) {
	t.Parallel()

	decorated := Decorate[string](Factory(constantBinder("plain")), suffixDecorator)

	views, ok := decorated.(elemViews)
	if !ok {
		t.Fatal("a decorated bind lost the elemType views, so duck typing cannot see it")
	}
	if got := views.ElementKey(); got != (types.KeyElem[string]{}) {
		t.Fatalf("ElementKey() = %#v, want KeyElem[string]", got)
	}

	generator, isGenerator := decorated.(types.AnyGenerator)
	if !isGenerator {
		t.Fatal("a decorated bind must still generate through the type-erased view")
	}
	erased, err := generator.GenAsAny(nil)
	if err != nil {
		t.Fatalf("GenAsAny() error = %v, want nil", err)
	}
	if text, isText := erased.(string); !isText || text != "plain-decorated" {
		t.Fatalf("GenAsAny() = %#v, want the decorated value", erased)
	}
}

func BenchmarkDecorate_Generates(b *testing.B) {
	decorated := Decorate[int](Factory(constantBinder(42)), func(
		_ types.DependencyRetriever, value int,
	) (int, error) {
		return value + 1, nil
	})

	var sink int
	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		sink, _ = decorated.Generates(nil)
	}
	runtime.KeepAlive(sink)
}

func BenchmarkFactoryBind_GeneratesUndecorated(b *testing.B) {
	plain := Factory(constantBinder(42))

	var sink int
	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		sink, _ = plain.Generates(nil)
	}
	runtime.KeepAlive(sink)
}

type foreignBind struct {
	value string
}

func (b foreignBind) Generates(types.DependencyRetriever) (string, error) {
	return b.value, nil
}

func (b foreignBind) Type() types.BindType {
	return types.BindFactory
}

func TestDecorate_WrapsABindItCannotRebuild(t *testing.T) {
	t.Parallel()

	decorated := Decorate[string](foreignBind{value: "plain"}, suffixDecorator)
	if _, wrapped := decorated.(DecoratedBind[string]); !wrapped {
		t.Fatalf("Decorate() = %T, want a bind it does not own to fall back to the wrapper", decorated)
	}

	got, err := decorated.Generates(nil)
	if err != nil {
		t.Fatalf("Generates() error = %v, want nil", err)
	}
	if got != "plain-decorated" {
		t.Fatalf("Generates() = %q, want the wrapped value decorated", got)
	}
	if kind := decorated.Type(); kind != types.BindFactory {
		t.Fatalf("Type() = %d, want the kind the wrapped bind reports", kind)
	}
}
