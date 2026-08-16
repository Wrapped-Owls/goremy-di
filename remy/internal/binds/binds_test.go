package binds

import (
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

var errBinderFailed = errors.New("binder failed")

func constantBinder[T any](value T) types.Binder[T] {
	return func(types.DependencyRetriever) (T, error) { return value, nil }
}

func failingBinder[T any]() types.Binder[T] {
	return func(types.DependencyRetriever) (result T, err error) { return result, errBinderFailed }
}

func countingBinder(calls *int) types.Binder[int] {
	return func(types.DependencyRetriever) (int, error) {
		*calls++
		return *calls, nil
	}
}

func TestBindConstructors_Type(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		bind types.Bind[int]
		want types.BindType
	}{
		{name: "Instance", bind: Instance(7), want: types.BindInstance},
		{name: "Factory", bind: Factory(constantBinder(7)), want: types.BindFactory},
		{name: "Singleton", bind: Singleton(constantBinder(7)), want: types.BindSingleton},
		{
			name: "LazySingleton",
			bind: LazySingleton(constantBinder(7)),
			want: types.BindLazySingleton,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := testCase.bind.Type(); got != testCase.want {
				t.Fatalf("Type() = %d, want %d", got, testCase.want)
			}
		})
	}
}

func TestBindConstructors_Generates(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		bind types.Bind[int]
	}{
		{name: "Instance", bind: Instance(7)},
		{name: "Factory", bind: Factory(constantBinder(7))},
		{name: "Singleton", bind: Singleton(constantBinder(7))},
		{name: "LazySingleton", bind: LazySingleton(constantBinder(7))},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.bind.Generates(nil)
			if err != nil {
				t.Fatalf("Generates() error = %v, want nil", err)
			}
			if got != 7 {
				t.Fatalf("Generates() = %d, want 7", got)
			}
		})
	}
}

func TestBindConstructors_GeneratesError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		bind types.Bind[int]
	}{
		{name: "Factory", bind: Factory(failingBinder[int]())},
		{name: "Singleton", bind: Singleton(failingBinder[int]())},
		{name: "LazySingleton", bind: LazySingleton(failingBinder[int]())},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if _, err := testCase.bind.Generates(nil); !errors.Is(err, errBinderFailed) {
				t.Fatalf("Generates() error = %v, want %v", err, errBinderFailed)
			}
		})
	}
}
