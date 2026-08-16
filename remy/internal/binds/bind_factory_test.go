package binds

import (
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

func TestFactoryBind_GeneratesEveryCall(t *testing.T) {
	t.Parallel()

	var calls int
	bind := Factory(countingBinder(&calls))

	for round := 1; round <= 3; round++ {
		got, err := bind.Generates(nil)
		if err != nil {
			t.Fatalf("call %d: unexpected error: %v", round, err)
		}
		if got != round {
			t.Fatalf("call %d returned %d, want a freshly built %d", round, got, round)
		}
	}
	if calls != 3 {
		t.Fatalf("binder ran %d times, want one run per retrieval", calls)
	}
}

func TestInstanceBind_GeneratesTheSameValue(t *testing.T) {
	t.Parallel()

	bind := Instance(7)
	for round := 1; round <= 3; round++ {
		got, err := bind.Generates(nil)
		if err != nil {
			t.Fatalf("call %d: unexpected error: %v", round, err)
		}
		if got != 7 {
			t.Fatalf("call %d returned %d, want the captured 7", round, got)
		}
	}
}

func TestFactoryBind_GenAsAny(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		bind    types.Bind[string]
		want    any
		wantErr error
	}{
		{
			name: "factory erases the built value",
			bind: Factory(constantBinder("gopher")),
			want: "gopher",
		},
		{
			name: "instance erases the captured value",
			bind: Instance("gopher"),
			want: "gopher",
		},
		{
			name:    "a failing binder surfaces its error",
			bind:    Factory(failingBinder[string]()),
			want:    "",
			wantErr: errBinderFailed,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			generator, ok := testCase.bind.(types.AnyGenerator)
			if !ok {
				t.Fatalf("%T does not implement types.AnyGenerator", testCase.bind)
			}

			got, err := generator.GenAsAny(nil)
			if !errors.Is(err, testCase.wantErr) {
				t.Fatalf("GenAsAny() error = %v, want %v", err, testCase.wantErr)
			}
			if got != testCase.want {
				t.Fatalf("GenAsAny() = %#v, want %#v", got, testCase.want)
			}
		})
	}
}
