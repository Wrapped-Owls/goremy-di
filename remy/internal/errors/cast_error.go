package errors

import (
	"fmt"
)

// ErrTypeCastInRuntime indicates that the library was unable to find/cast the element with the given type.
type ErrTypeCastInRuntime struct {
	baseErrorChecker[ErrTypeCastInRuntime, *ErrTypeCastInRuntime]
	ActualValue any
	Expected    any
}

func (e ErrTypeCastInRuntime) Error() string {
	if e.ActualValue == nil {
		return fmt.Sprintf("unable to find/cast the element with given type: %T", e.Expected)
	}
	return fmt.Sprintf("unable to cast `%T` to given type `%T`", e.ActualValue, e.Expected)
}
