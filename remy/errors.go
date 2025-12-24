package remy

import (
	"github.com/wrapped-owls/goremy-di/remy/internal/errors"
)

// Re-export errors from internal/errors package for backward compatibility
var (
	ErrAlreadyBound            = errors.ErrAlreadyBoundSentinel
	ErrImpossibleIdentifyType  = errors.ErrImpossibleIdentifyTypeSentinel
	ErrElementNotRegistered    = errors.ErrElementNotRegisteredSentinel
	ErrCycleDependencyDetected = errors.ErrCycleDependencyDetectedSentinel
	ErrTypeCastInRuntime       = errors.ErrTypeCastInRuntimeSentinel
	ErrFoundMoreThanOneValidDI = errors.ErrFoundMoreThanOneValidDISentinel
)
