package types

type (
	CheckReflectionOptions interface {
		ReflectOpts() ReflectionOptions
	}
	ValuesSetter[K comparable, T any] interface {
		// Set uses the T key given to save the value.
		// If the key is already bound, it returns a boolean with value true.
		// Returns an error if the key cannot be set (e.g., override not allowed).
		Set(K, T) (wasOverridden bool, err error)

		// SetNamed uses the T key and cacheKey given to store the value.
		// If the key is already bound, it returns a boolean with value true.
		// Returns an error if the key cannot be set (e.g., override not allowed).
		SetNamed(K, string, T) (wasOverridden bool, err error)
		CheckReflectionOptions
	}
	AllValuesGetter[T any] interface {
		// GetAll returns all elements from the storage that hasn't a key
		GetAll(optKey ...string) ([]T, error)
	}
	ValuesGetter[K comparable, T any] interface {
		// GetNamed search for a named element that was cached using the T value given and a string key
		GetNamed(K, string) (T, error)

		// Get search for a named element that was cached using the T value given
		Get(K) (T, error)

		AllValuesGetter[T]
		CheckReflectionOptions
	}

	// Storage is the main cache interface that is used by the injector to store the values
	Storage[K comparable] interface {
		ValuesSetter[K, any]
		ValuesGetter[K, any]
	}

	// DependencyRetriever is the main element used to obtain registered binds/instances
	DependencyRetriever interface {
		WrapRetriever() Injector
		RetrieveBind(bindKey BindKey, tag string) (any, error)

		AllValuesGetter[any]
		CheckReflectionOptions
	}

	// Injector is the main interface that contains all needed methods to make an injector work
	Injector interface {
		BindElem(depKey BindKey, val any, opts BindOptions) error
		SubInjector(allowOverrides ...bool) Injector
		DependencyRetriever
	}
)
