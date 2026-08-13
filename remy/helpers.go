package remy

import (
	"errors"

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

// every constructor builds its base the same way, so none of them can forget to
// honour ParentInjector
func newBaseInjector(conf Config) types.Injector {
	opts := injectorOptsFromConfig(conf)
	if conf.ParentInjector != nil {
		return injector.New(opts, conf.ParentInjector)
	}

	return injector.New(opts)
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

// split out because recover only works when called directly by the deferred func
func applyRecovered(recovered any, err *error) {
	asError := remyErrs.FromRecovered(recovered)
	if asError == nil || err == nil {
		return
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
