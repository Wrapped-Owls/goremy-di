package errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

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
	ErrDependencyResolutionSentinel    = &ErrDependencyResolution{}
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

// FromRecovered turns a recovered panic value into an error, nil when nothing
// was recovered. Callers must pass the result of recover directly.
func FromRecovered(recovered any) error {
	if recovered == nil {
		return nil
	}
	if asError, ok := recovered.(error); ok {
		return asError
	}

	return fmt.Errorf("%v", recovered)
}

func debugBindKey(value types.BindKey) (keyVal string) {
	if value == nil {
		return ""
	}

	keyVal = reflect.TypeOf(value).String()
	return " " + keyVal
}

func writePathEntry(builder *strings.Builder, entry PathEntry) {
	builder.WriteString(debugBindKey(entry.Key))
	if entry.Tag != "" {
		builder.WriteString(` (tag "`)
		builder.WriteString(entry.Tag)
		builder.WriteString(`")`)
	}
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

func CheckError[T any](checkErr any) (foundErr T, ok bool) {
	type errUnwrap interface {
		Unwrap() error
	}

	for {
		asUnwrap, canUnwrap := checkErr.(errUnwrap)
		if !canUnwrap || asUnwrap == nil {
			break
		}

		unwrapped := asUnwrap.Unwrap()
		if unwrapped == nil {
			checkErr = nil
			break
		}

		checkErr = unwrapped
	}

	switch val := checkErr.(type) {
	case T:
		foundErr = val
	case *T:
		if val == nil {
			return foundErr, false
		}
		foundErr = *val
	default:
		return foundErr, false
	}

	return foundErr, true
}
