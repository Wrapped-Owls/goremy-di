package utils

import "errors"

var (
	ErrAlreadyBound                 = errors.New("dependency already bound")
	ErrImpossibleIdentifyType       = errors.New("impossible to identify the given type")
	ErrElementNotRegistered         = errors.New("element with given key is not registered")
	ErrNoElementFoundInsideOrParent = errors.New("no element found on the given injector or any of it's parents")
	ErrCycleDependencyDetected      = errors.New("cycle dependency detected, check for it")
)
