package graph

import (
	"sync"
	"testing"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

func node(name string, key types.BindKey) Node {
	return Node{Key: key, Tag: name}
}

func TestState_RecordEdge(t *testing.T) {
	t.Parallel()

	first := Edge{From: node("a", types.KeyElem[string]{}), To: node("b", types.KeyElem[int]{})}
	second := Edge{From: node("b", types.KeyElem[int]{}), To: node("c", types.KeyElem[bool]{})}
	sameAsFirstOtherTag := Edge{
		From: node("z", types.KeyElem[string]{}),
		To:   node("b", types.KeyElem[int]{}),
	}

	testCases := []struct {
		name      string
		record    []Edge
		wantEdges []Edge
	}{
		{
			name:      "records each distinct edge once",
			record:    []Edge{first, second},
			wantEdges: []Edge{first, second},
		},
		{
			name:      "repeating an edge does not duplicate it",
			record:    []Edge{first, second, first, first},
			wantEdges: []Edge{first, second},
		},
		{
			name:      "same target from a different tag is a distinct edge",
			record:    []Edge{first, sameAsFirstOtherTag},
			wantEdges: []Edge{first, sameAsFirstOtherTag},
		},
		{
			name:      "nothing recorded leaves no edge",
			record:    nil,
			wantEdges: nil,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase // go.mod pins 1.20: loop vars are still shared
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			tracked := newState()
			for _, edge := range testCase.record {
				tracked.recordEdge(edge)
			}

			if len(tracked.edges) != len(testCase.wantEdges) {
				t.Fatalf("edges = %v, want %v", tracked.edges, testCase.wantEdges)
			}
			for index, want := range testCase.wantEdges {
				if tracked.edges[index] != want {
					t.Errorf("edges[%d] = %v, want %v", index, tracked.edges[index], want)
				}
			}
		})
	}
}

// the mutex exists because independent Get chains record concurrently
func TestState_RecordEdgeIsRaceFree(t *testing.T) {
	t.Parallel()

	const writers = 8
	tracked := newState()
	shared := Edge{From: node("a", types.KeyElem[string]{}), To: node("b", types.KeyElem[int]{})}

	var waitGroup sync.WaitGroup
	start := make(chan struct{})
	waitGroup.Add(writers)
	for writer := 0; writer < writers; writer++ {
		go func(writer int) {
			defer waitGroup.Done()
			<-start // release every goroutine at once so they genuinely contend

			tracked.recordEdge(shared)
			tracked.recordEdge(Edge{
				From: shared.From,
				To:   node(string(rune('a'+writer)), types.KeyElem[bool]{}),
			})
		}(writer)
	}
	close(start)
	waitGroup.Wait()

	// the shared edge is deduped to one, each writer adds its own distinct edge
	if len(tracked.edges) != writers+1 {
		t.Fatalf("edges = %d, want %d", len(tracked.edges), writers+1)
	}
	if snapshot := tracked.snapshotEdges(); len(snapshot) != writers+1 {
		t.Fatalf("snapshotEdges = %d, want %d", len(snapshot), writers+1)
	}
}
