package errors

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

// ErrCycleDependencyDetected indicates that a cycle dependency was detected.
type ErrCycleDependencyDetected struct {
	baseErrorChecker[ErrCycleDependencyDetected, *ErrCycleDependencyDetected]
	Path types.DependencyGraph
}

func (e *ErrCycleDependencyDetected) Error() string {
	return "cycle dependency detected, check for it"
}
