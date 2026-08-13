package remy

import (
	"errors"
	"fmt"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/injector"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

func injectorOptsFromConfig(conf Config) injector.Options {
	var opts injector.Options

	if conf.CanOverride {
		opts.Cache |= injopts.CacheOptAllowOverride
	}

	if conf.DuckTypeElements {
		opts.Cache |= injopts.CacheOptReturnAll
		opts.Resolve |= injopts.ResolveOptDuckTyping
	}

	if conf.MultiBinding {
		opts.Cache |= injopts.CacheOptReturnAll
	}

	return opts
}

func recoverInjectorPanic(err *error) {
	applyRecovered(recover(), err)
}

// deferred as a single direct call, never a closure, so recover() still works and
// the success path keeps costing one open-coded defer
func traceResolution[T any](err *error, retriever DependencyRetriever, tag string) {
	applyRecovered(recover(), err)
	if *err == nil {
		return
	}

	*err = remyErrs.WrapResolutionPath(*err, utils.NewKeyElem[T](), tag)
}

// separate from recoverInjectorPanic because recover() only works when called
// directly by the deferred function
func applyRecovered(r any, err *error) {
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
