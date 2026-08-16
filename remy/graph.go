package remy

import "github.com/wrapped-owls/goremy-di/remy/internal/graph"

// NewCycleDetector creates an Injector that detects cycle dependencies at runtime
// without recording any edge. Being much slower than the standard one, it suits tests.
//
//nolint:ireturn // DI container constructor: the decorator is an implementation detail
func NewCycleDetector(configs ...Config) Injector {
	return graph.NewDetector(newBaseInjector(firstOrDefault(configs...)))
}

// NewCycleDetectorInjector creates an Injector that detects cycle dependencies.
//
// Deprecated: use NewCycleDetector.
//
//go:fix inline
//nolint:ireturn // DI container constructor: kept for backwards compatibility
func NewCycleDetectorInjector(configs ...Config) Injector {
	return NewCycleDetector(configs...)
}
