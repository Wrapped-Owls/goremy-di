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
)
