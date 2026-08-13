package types

import "github.com/wrapped-owls/goremy-di/remy/pkg/injopts"

type (
	ValuesSetter[K comparable, T any] interface {
		// Set uses the T key given to save the value.
		// If the key is already bound, it returns a boolean with value true.
		// Returns an error if the key cannot be set (e.g., override not allowed).
		Set(K, T) (wasOverridden bool, err error)

		// SetNamed uses the T key and cacheKey given to store the value.
		// If the key is already bound, it returns a boolean with value true.
		// Returns an error if the key cannot be set (e.g., override not allowed).
		SetNamed(K, string, T) (wasOverridden bool, err error)
	}
	AllValuesGetter[T any] interface {
		// GetAll returns all elements from the storage that hasn't a key
		GetAll(keyTag string) ([]T, error)
	}
	ValuesGetter[K comparable, T any] interface {
		// GetNamed search for a named element that was cached using the T value given and a string key
		GetNamed(K, string) (T, error)

		// Get search for a named element that was cached using the T value given
		Get(K) (T, error)

		AllValuesGetter[T]
	}

	// Storage is the main cache interface that is used by the injector to store the values
	Storage[K comparable] interface {
		ValuesSetter[K, any]
		ValuesGetter[K, any]
	}

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

	// AnyGenerator generates an element without knowing its type parameter
	AnyGenerator interface {
		GenAsAny(retriever DependencyRetriever) (any, error)
	}

	// GuessableBind tells whether a bind satisfies an interface before generating it
	GuessableBind interface {
		PointerValue() any
		DefaultValue() any
		AnyGenerator
	}
)
