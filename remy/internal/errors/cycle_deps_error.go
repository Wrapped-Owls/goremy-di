package errors

import "strings"

// ErrCycleDependencyDetected indicates that a cycle dependency was detected.
// Path holds the resolution chain outermost-first, ending on the node that
// closed the cycle.
type ErrCycleDependencyDetected struct {
	baseErrorChecker[ErrCycleDependencyDetected, *ErrCycleDependencyDetected]
	Path []PathEntry
}

func (e ErrCycleDependencyDetected) Error() string {
	if len(e.Path) == 0 {
		return "cycle dependency detected, check for it"
	}

	var builder strings.Builder
	builder.WriteString("cycle dependency detected:")
	for index, node := range e.Path {
		if index > 0 {
			builder.WriteString(" ->")
		}
		writePathEntry(&builder, node)
	}
	return builder.String()
}
