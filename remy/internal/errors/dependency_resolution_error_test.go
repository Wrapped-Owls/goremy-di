package errors

import (
	"errors"
	"strings"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

func TestWrapResolutionPath(t *testing.T) {
	t.Parallel()

	rootCause := ErrElementNotRegisteredSentinel

	testCases := []struct {
		name         string
		buildErr     func() error
		wantPathLen  int
		wantContains []string
	}{
		{
			name: "single level creates path with one entry",
			buildErr: func() error {
				return WrapResolutionPath(rootCause, types.KeyElem[bool]{}, "")
			},
			wantPathLen:  1,
			wantContains: []string{"bool"},
		},
		{
			name: "nested levels append in place keeping one error",
			buildErr: func() error {
				err := WrapResolutionPath(rootCause, types.KeyElem[bool]{}, "")
				err = WrapResolutionPath(err, types.KeyElem[int]{}, "")
				return WrapResolutionPath(err, types.KeyElem[string]{}, "")
			},
			wantPathLen:  3,
			wantContains: []string{"string", "int", "bool"},
		},
		{
			name: "tagged entry renders its tag",
			buildErr: func() error {
				return WrapResolutionPath(rootCause, types.KeyElem[int]{}, "primary")
			},
			wantPathLen:  1,
			wantContains: []string{"int", `(tag "primary")`},
		},
		{
			name: "chain deeper than the inline storage spills to overflow",
			buildErr: func() error {
				err := WrapResolutionPath(rootCause, types.KeyElem[bool]{}, "")
				err = WrapResolutionPath(err, types.KeyElem[int8]{}, "")
				err = WrapResolutionPath(err, types.KeyElem[int16]{}, "")
				err = WrapResolutionPath(err, types.KeyElem[int32]{}, "")
				err = WrapResolutionPath(err, types.KeyElem[int64]{}, "")
				return WrapResolutionPath(err, types.KeyElem[string]{}, "")
			},
			wantPathLen:  6,
			wantContains: []string{"string", "int64", "int32", "int16", "int8", "bool"},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := testCase.buildErr()
			depErr, ok := err.(*ErrDependencyResolution)
			if !ok {
				t.Fatalf("expected *ErrDependencyResolution, got %T", err)
			}
			if pathLen := len(depErr.Path()); pathLen != testCase.wantPathLen {
				t.Fatalf("path length = %d, want %d", pathLen, testCase.wantPathLen)
			}
			if !errors.Is(err, ErrElementNotRegisteredSentinel) {
				t.Errorf("root cause lost through wrapping: %v", err)
			}

			errMsg := err.Error()
			for _, fragment := range testCase.wantContains {
				if !strings.Contains(errMsg, fragment) {
					t.Errorf("message %q does not contain %q", errMsg, fragment)
				}
			}
		})
	}
}

func TestWrapResolutionPath_RendersOutermostFirst(t *testing.T) {
	t.Parallel()

	err := WrapResolutionPath(ErrElementNotRegisteredSentinel, types.KeyElem[bool]{}, "")
	err = WrapResolutionPath(err, types.KeyElem[string]{}, "")

	errMsg := err.Error()
	outerIndex := strings.Index(errMsg, "string")
	innerIndex := strings.Index(errMsg, "bool")
	if outerIndex < 0 || innerIndex < 0 || outerIndex > innerIndex {
		t.Fatalf("expected outermost type first in %q", errMsg)
	}
}

// note: AllocsPerRun forbids parallel tests, so this one must not call t.Parallel
func TestWrapResolutionPath_ConstantAllocations(t *testing.T) {
	// depth 3 must cost the same single allocation as depth 1: the error
	// itself, whose inline array stores the whole chain
	allocs := testing.AllocsPerRun(100, func() {
		err := WrapResolutionPath(ErrElementNotRegisteredSentinel, types.KeyElem[bool]{}, "")
		err = WrapResolutionPath(err, types.KeyElem[int]{}, "")
		err = WrapResolutionPath(err, types.KeyElem[string]{}, "")
		if err == nil {
			t.Fatal("expected error")
		}
	})
	if allocs > 1 {
		t.Fatalf("wrapping a depth-3 chain allocated %.0f times, want 1", allocs)
	}
}
