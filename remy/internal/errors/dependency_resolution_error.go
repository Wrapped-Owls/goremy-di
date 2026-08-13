package errors

import (
	"strings"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

// PathEntry identifies one step of a dependency-resolution chain
type PathEntry = types.PathEntry

// keeps the common A -> B -> C chain at a single allocation; deeper chains spill
const inlinePathDepth = 3

// ErrDependencyResolution carries the resolution path that led to Cause, stored
// innermost-first. Outer levels append in place instead of stacking one wrapper
// error each, so a whole chain costs a single allocation.
type ErrDependencyResolution struct {
	baseErrorChecker[ErrDependencyResolution, *ErrDependencyResolution]
	Cause    error
	entries  [inlinePathDepth]PathEntry
	depth    int
	overflow []PathEntry
}

// WrapResolutionPath records (key, tag) as one more resolution level on top of
// cause, appending in place when cause already carries a path.
func WrapResolutionPath(cause error, key types.BindKey, tag string) error {
	// direct assertion, not errors.As: nested levels always get it unwrapped
	if depErr, ok := cause.(*ErrDependencyResolution); ok {
		depErr.push(key, tag)
		return depErr
	}

	depErr := &ErrDependencyResolution{Cause: cause}
	depErr.push(key, tag)
	return depErr
}

func (e *ErrDependencyResolution) push(key types.BindKey, tag string) {
	if e.depth < inlinePathDepth {
		e.entries[e.depth] = PathEntry{Key: key, Tag: tag}
	} else {
		e.overflow = append(e.overflow, PathEntry{Key: key, Tag: tag})
	}
	e.depth++
}

// Path returns the resolution chain innermost-first.
func (e *ErrDependencyResolution) Path() []PathEntry {
	inline := e.depth
	if inline > inlinePathDepth {
		inline = inlinePathDepth
	}
	path := make([]PathEntry, 0, e.depth)
	path = append(path, e.entries[:inline]...)
	return append(path, e.overflow...)
}

func (e *ErrDependencyResolution) Unwrap() error {
	return e.Cause
}

func (e *ErrDependencyResolution) Error() string {
	var builder strings.Builder
	builder.WriteString("failed to resolve")

	// stored innermost-first; render outermost-first so it reads A -> B -> C
	path := e.Path()
	for index := len(path) - 1; index >= 0; index-- {
		entry := path[index]
		builder.WriteString(debugBindKey(entry.Key))
		if entry.Tag != "" {
			builder.WriteString(` (tag "`)
			builder.WriteString(entry.Tag)
			builder.WriteString(`")`)
		}
		if index > 0 {
			builder.WriteString(" ->")
		}
	}

	builder.WriteString(": ")
	if e.Cause != nil {
		builder.WriteString(e.Cause.Error())
	}
	return builder.String()
}
