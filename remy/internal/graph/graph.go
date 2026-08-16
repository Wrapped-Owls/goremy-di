package graph

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

type (
	Node = types.GraphNode
	Edge = types.GraphEdge

	FailedNode struct {
		Node Node
		Err  error
	}

	Graph interface {
		// Edges returns every dependent -> dependency pair recorded so far
		Edges() []Edge

		// ResolveAll generates every registered bind to record its edges, and
		// reports the ones that failed instead of aborting on the first
		ResolveAll() ([]FailedNode, error)
	}
)
