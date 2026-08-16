package remy

import "github.com/wrapped-owls/goremy-di/remy/internal/graph"

type (
	// GraphNode identifies one dependency in a resolution graph
	GraphNode = graph.Node

	// GraphEdge records that From requested To during a resolution
	GraphEdge = graph.Edge

	// FailedNode is a bind ResolveAll could not build, with the reason it gave
	FailedNode = graph.FailedNode

	// Graph exposes the dependency graph recorded by a graph injector
	Graph = graph.Graph
)

// NewGraphInjector creates an Injector that records the dependency graph as elements
// resolve, plus the Graph view to read it. It suits debugging, never a hot path.
//
//nolint:ireturn // DI container constructor: both views are interfaces by design
func NewGraphInjector(configs ...Config) (Injector, Graph) {
	return graph.New(newBaseInjector(firstOrDefault(configs...)))
}

// NewCycleDetector creates an Injector that detects cycle dependencies at runtime
// without recording any edge. Being much slower than the standard one, it suits tests.
//
//nolint:ireturn // DI container constructor: the decorator is an implementation detail
func NewCycleDetector(configs ...Config) Injector {
	return graph.NewDetector(newBaseInjector(firstOrDefault(configs...)))
}

// NewCycleDetectorInjector creates an Injector that detects cycle dependencies.
//
// Deprecated: use NewCycleDetector, or NewGraphInjector when the recorded
// dependency graph is also wanted.
//
//go:fix inline
//nolint:ireturn // DI container constructor: kept for backwards compatibility
func NewCycleDetectorInjector(configs ...Config) Injector {
	return NewCycleDetector(configs...)
}
