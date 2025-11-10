package errors

import "fmt"

// ErrWrapParentSubErrors indicates that no element was found in the given injector or any of its parents.
type ErrWrapParentSubErrors struct {
	baseErrorChecker[ErrWrapParentSubErrors, *ErrWrapParentSubErrors]
	MainError error
	SubError  error
}

func (e ErrWrapParentSubErrors) Unwrap() []error {
	return []error{e.MainError, e.SubError}
}

func (e ErrWrapParentSubErrors) Error() string {
	return fmt.Sprintf("error in parent injector, main error is: %v", e.MainError)
}
