package errors

import (
	"fmt"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

// ErrAlreadyBound indicates that a dependency has already been bound and cannot be registered again.
type ErrAlreadyBound struct {
	baseErrorChecker[ErrAlreadyBound, *ErrAlreadyBound]
	Key types.BindKey
}

func (e ErrAlreadyBound) Error() string {
	debugKey := debugBindKey(e.Key)
	return fmt.Sprintf("dependency already bound: %s", debugKey)
}

// ErrElementNotRegistered indicates that an element with the given key is not registered.
type ErrElementNotRegistered struct {
	baseErrorChecker[ErrElementNotRegistered, *ErrElementNotRegistered]
	Key any
}

func (e ErrElementNotRegistered) Error() string {
	debugKey := ""
	if bindKey, ok := e.Key.(types.BindKey); ok {
		debugKey = debugBindKey(bindKey)
	} else {
		debugKey = genDebugKeyTypeName(e.Key)
	}
	return fmt.Sprintf("element with given key (`%s`) is not registered", debugKey)
}
