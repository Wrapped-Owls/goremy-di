package errors

import (
	"errors"
	"reflect"

	"github.com/wrapped-owls/goremy-di/remy/internal/types"
)

// Sentinel errors for backward compatibility and easier error checking
var (
	ErrAlreadyBoundSentinel           = &ErrAlreadyBound{}
	ErrImpossibleIdentifyTypeSentinel = &ErrImpossibleIdentifyType{}
	ErrElementNotRegisteredSentinel   = &ErrElementNotRegistered{}
	ErrConfigNotAllowReturnAll        = errors.New(
		"the current injector config does not allow returning all elements",
	)
	ErrCycleDependencyDetectedSentinel = &ErrCycleDependencyDetected{}
	ErrTypeCastInRuntimeSentinel       = &ErrTypeCastInRuntime{}
	ErrFoundMoreThanOneValidDISentinel = &ErrMultipleDIDuckTypingCandidates{}
)

func genDebugKeyTypeName(typeKey any) (givenType string) {
	if typeKey != nil {
		if asReflectVal, ok := typeKey.(reflect.Type); ok && asReflectVal != nil {
			givenType = asReflectVal.Name()
		} else {
			givenType = reflect.TypeOf(typeKey).Name()
		}

		givenType = " `" + givenType + "`"
	}
	return givenType
}

func debugBindKey(value types.BindKey) (keyVal string) {
	if value == nil {
		return ""
	}

	keyVal = reflect.TypeOf(value).String()
	return " " + keyVal
}

type errorInterface[T any] interface {
	*T
	error
}

type baseErrorChecker[T any, PT errorInterface[T]] struct{}

func (e baseErrorChecker[T, PT]) Is(target error) bool {
	var asPointer PT
	if errors.As(target, &asPointer) {
		return true
	}

	_, ok := target.(T) // Check the raw value directly
	return ok
}
