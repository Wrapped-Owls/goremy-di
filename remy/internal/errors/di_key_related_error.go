package errors

import (
	"fmt"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

type (
	messageKeyRelatedGen interface {
		genErrorFormat(debugKey string) string
	}
	errDIKeyRelated[T messageKeyRelatedGen] struct {
		baseErrorChecker[errDIKeyRelated[T], *errDIKeyRelated[T]]
		Key any
	}
)

func (e errDIKeyRelated[T]) Error() string {
	debugKey := ""
	if bindKey, ok := e.Key.(types.BindKey); ok {
		debugKey = debugBindKey(bindKey)
	} else {
		debugKey = genDebugKeyTypeName(e.Key)
	}

	var formatGen T
	return formatGen.genErrorFormat(debugKey)
}

type (
	keyAlreadyBoundMessageGen struct{}
	// ErrAlreadyBound indicates that a dependency has already been bound and cannot be registered again.
	ErrAlreadyBound = errDIKeyRelated[keyAlreadyBoundMessageGen]
)

func (e keyAlreadyBoundMessageGen) genErrorFormat(debugKey string) string {
	return fmt.Sprintf("dependency already bound: %s", debugKey)
}

type (
	keyUnregisteredMessageGen struct{}
	// ErrElementNotRegistered indicates that an element with the given key is not registered.
	ErrElementNotRegistered = errDIKeyRelated[keyUnregisteredMessageGen]
)

func (e keyUnregisteredMessageGen) genErrorFormat(debugKey string) string {
	return fmt.Sprintf("element with given key (`%s`) is not registered", debugKey)
}
