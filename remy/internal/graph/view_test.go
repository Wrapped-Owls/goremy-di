package graph

import (
	"errors"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/binds"
	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/injcontainer"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

func TestView_ResolveAll(t *testing.T) {
	t.Parallel()

	inj, graph := New(newBase())
	registerChain(t, inj)

	// this one depends on a type nobody registered
	if err := injcontainer.Register[float32](inj, "", binds.Factory(
		func(r types.DependencyRetriever) (float32, error) {
			_, err := injcontainer.Get[uint8](r, "")
			return 0, err
		},
	)); err != nil {
		t.Fatalf("register orphan: %v", err)
	}

	failures, err := graph.ResolveAll()
	if !errors.Is(err, remyErrs.ErrElementNotRegisteredSentinel) {
		t.Fatalf("expected joined not-registered error, got: %v", err)
	}

	if len(failures) != 1 || failures[0].Node.Key != (types.KeyElem[float32]{}) {
		t.Fatalf("failures = %v, want only float32", failures)
	}
	if !errors.Is(failures[0].Err, remyErrs.ErrElementNotRegisteredSentinel) {
		t.Fatalf("failures[0].Err = %v, want the reason it failed with", failures[0].Err)
	}

	// forcing resolution recorded the lazy chain even without a single Get
	if edges := graph.Edges(); len(edges) < 3 {
		t.Fatalf("edges = %d, want at least 3: %v", len(edges), edges)
	}
}

func TestView_EdgesAreASnapshot(t *testing.T) {
	t.Parallel()

	inj, graph := New(newBase())
	registerChain(t, inj)

	if _, err := injcontainer.Get[string](inj, ""); err != nil {
		t.Fatalf("unexpected resolution error: %v", err)
	}

	edges := graph.Edges()
	total := len(edges)
	edges[0] = Edge{}

	if again := graph.Edges(); len(again) != total || again[0] == (Edge{}) {
		t.Fatal("Edges must hand out a copy, not the recorded slice")
	}
}
