package remy

import (
	"errors"
	"fmt"

	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
)

func cacheOptsFromConfig(conf Config) (options injopts.CacheConfOption) {
	if conf.CanOverride {
		options |= injopts.CacheOptAllowOverride
	}

	if conf.DuckTypeElements {
		options |= injopts.CacheOptReturnAll
	}

	return
}

func recoverInjectorPanic(err *error) {
	r := recover()
	if r == nil || err == nil {
		return
	}

	var asError error
	switch asVal := r.(type) {
	case error:
		asError = asVal
	default:
		asError = fmt.Errorf("%v", r)
	}

	if *err != nil {
		*err = errors.Join(*err, asError)
		return
	}

	*err = asError
}

func firstOrDefault[T any](valueList ...T) (result T) {
	if len(valueList) > 0 {
		result = valueList[0]
	}

	return result
}
