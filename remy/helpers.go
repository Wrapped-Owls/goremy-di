package remy

import (
	"errors"
	"fmt"

	remyErrs "github.com/wrapped-owls/goremy-di/remy/internal/errors"
	"github.com/wrapped-owls/goremy-di/remy/internal/injector"
	"github.com/wrapped-owls/goremy-di/remy/internal/types"
	"github.com/wrapped-owls/goremy-di/remy/pkg/injopts"
	"github.com/wrapped-owls/goremy-di/remy/pkg/utils"
)

func injectorOptsFromConfig(conf Config) injector.Options {
	opts := injector.Options{ScopeName: conf.ScopeName}

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

	if conf.TraceResolution {
		opts.Resolve |= injopts.ResolveOptTracePath
	}

	if conf.Isolated {
		opts.Resolve |= injopts.ResolveOptIsolated
	}

	return opts
}

func recoverInjectorPanic(err *error) {
	applyRecovered(recover(), err)
}

// traceResolution recovers an injector panic and records key as one more step of
// the failing path. It must be deferred directly, never in a closure, or recover
// stops working.
func traceResolution[T any](err *error, retriever DependencyRetriever, tag string) {
	applyRecovered(recover(), err)
	if *err == nil {
		return
	}

	if holder, ok := retriever.(types.ResolveOptionsHolder); ok &&
		holder.ResolveOptions().Is(injopts.ResolveOptTracePath) {
		*err = remyErrs.WrapResolutionPath(*err, utils.NewKeyElem[T](), tag)
	}
}

// applyRecovered folds a recovered panic value into err. It is split out because
// recover only works when called directly by the deferred function.
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
