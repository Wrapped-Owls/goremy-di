package binds

import (
	"errors"
	"runtime"
	"sync"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

func TestSingletonBind_GeneratesOnce(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		bind func(types.Binder[int]) types.Bind[int]
	}{
		{name: "Singleton", bind: Singleton[int]},
		{name: "LazySingleton", bind: LazySingleton[int]},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			const readers = 10
			var calls int
			bind := testCase.bind(countingBinder(&calls))

			var (
				waitGroup sync.WaitGroup
				results   [readers]int
			)
			start := make(chan struct{})
			waitGroup.Add(readers)
			for reader := 0; reader < readers; reader++ {
				go func(reader int) {
					defer waitGroup.Done()
					<-start // release every goroutine at once so they genuinely contend

					value, err := bind.Generates(nil)
					if err != nil {
						t.Errorf("reader %d: unexpected error: %v", reader, err)
					}
					results[reader] = value
				}(reader)
			}
			close(start)
			waitGroup.Wait()

			if calls != 1 {
				t.Fatalf("binder ran %d times, want exactly once", calls)
			}
			for reader, got := range results {
				if got != 1 {
					t.Fatalf(
						"reader %d saw %d, want every reader to share the one value",
						reader,
						got,
					)
				}
			}
		})
	}
}

// a failing binder must not be remembered as a successful build
func TestSingletonBind_BinderErrorIsNotCached(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		bind func(types.Binder[int]) types.Bind[int]
	}{
		{name: "Singleton", bind: Singleton[int]},
		{name: "LazySingleton", bind: LazySingleton[int]},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			bind := testCase.bind(failingBinder[int]())
			for round := 1; round <= 2; round++ {
				if _, err := bind.Generates(nil); !errors.Is(err, errBinderFailed) {
					t.Fatalf("call %d: error = %v, want %v", round, err, errBinderFailed)
				}
			}
		})
	}
}

func TestSingletonBind_RecoversAfterAFailedBuild(t *testing.T) {
	t.Parallel()

	var calls int
	bind := LazySingleton(func(types.DependencyRetriever) (int, error) {
		calls++
		if calls == 1 {
			return 0, errBinderFailed
		}
		return 42, nil
	})

	if _, err := bind.Generates(nil); !errors.Is(err, errBinderFailed) {
		t.Fatalf("first call error = %v, want %v", err, errBinderFailed)
	}

	got, err := bind.Generates(nil)
	if err != nil {
		t.Fatalf("second call error = %v, want nil", err)
	}
	if got != 42 {
		t.Fatalf("second call = %d, want the value the retry produced", got)
	}
}

func TestSingletonBind_GenAsAny(t *testing.T) {
	t.Parallel()

	generator, ok := Singleton(constantBinder("gopher")).(types.AnyGenerator)
	if !ok {
		t.Fatal("SingletonBind does not implement types.AnyGenerator")
	}

	got, err := generator.GenAsAny(nil)
	if err != nil {
		t.Fatalf("GenAsAny() error = %v, want nil", err)
	}
	if got != "gopher" {
		t.Fatalf("GenAsAny() = %#v, want the string gopher", got)
	}
}

// the already-built path is the hot one: a singleton builds once and is read forever
func BenchmarkSingletonBind_Generates(b *testing.B) {
	bind := Singleton(constantBinder(42))
	if _, err := bind.Generates(nil); err != nil {
		b.Fatalf("warm up: %v", err)
	}

	var sink int
	b.ReportAllocs()
	b.ResetTimer()
	for iteration := 0; iteration < b.N; iteration++ {
		sink, _ = bind.Generates(nil)
	}
	runtime.KeepAlive(sink)
}
