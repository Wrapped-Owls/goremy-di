package utils

import "errors"

var (
	ErrAlreadyBound           = errors.New("dependency already bound")
	ErrImpossibleIdentifyType = errors.New("impossible to identify the given type")
)
