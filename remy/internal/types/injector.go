package types

import "github.com/wrapped-owls/goremy-di/remy/pkg/injopts"

type (
	// DependencyRetriever is the main element used to obtain registered binds/instances
	DependencyRetriever interface {
		// RetrieverFor returns the retriever to use while resolving key and its
		// dependencies, or nil to keep using this one
		RetrieverFor(key BindKey, tag string) Injector

		RetrieveBind(bindKey BindKey, tag string) (any, error)

		AllValuesGetter[any]
	}

	// Injector is the main interface that contains all needed methods to make an injector work
	Injector interface {
		BindElem(depKey BindKey, val any, opts BindOptions) error
		SubInjector(allowOverrides ...bool) Injector

		// ScopeName returns the name given to this scope, empty when anonymous
		ScopeName() string

		// Parent returns the retriever this scope falls back to, nil on a root
		Parent() DependencyRetriever

		DependencyRetriever
	}

	// ResolveOptionsHolder exposes the resolution options in effect for an injector
	ResolveOptionsHolder interface {
		ResolveOptions() injopts.ResolveConfOption
	}

	// ScopeFactory builds the temporary scope a GetWith call resolves through
	ScopeFactory func(parent DependencyRetriever, storage Storage[BindKey]) Injector

	// BindEnumerator lists the binds an injector owns, excluding inherited ones
	BindEnumerator interface {
		ForEachBind(visit func(tag string, value any) bool)
	}

	// ScopeDecorator wraps a ScopeFactory so a decorated injector keeps tracking across
	// a temporary scope, where a cycle would otherwise recurse until the stack dies
	ScopeDecorator interface {
		WrapScopeFactory(base ScopeFactory) ScopeFactory
	}
)

// ResolveOptionsOf returns the options the given injector resolves with, and
// whether it exposes them at all.
func ResolveOptionsOf(value any) (opts injopts.ResolveConfOption, ok bool) {
	var holder ResolveOptionsHolder
	holder, ok = value.(ResolveOptionsHolder)
	if !ok {
		return injopts.ResolveOptNone, false
	}

	return holder.ResolveOptions(), true
}
