package types

type (
	CheckReflectionOptions interface {
		ReflectOpts() ReflectionOptions
	}
	ValuesSetter[T comparable] interface {
		// Set uses the T key given to save the value.
		// If the key is already bound, it returns a boolean with value true.
		Set(T, any) bool

		// SetNamed uses the T key and cacheKey given to store the value.
		// If the key is already bound, it returns a boolean with value true.
		SetNamed(T, string, any) bool
		CheckReflectionOptions
	}
	ValuesGetter[T comparable] interface {
		// GetNamed search for a named element that was cached using the T value given and a string key
		GetNamed(T, string) (any, error)

		// Get search for a named element that was cached using the T value given
		Get(T) (any, error)

		// GetAll returns all elements from the storage that hasn't a key
		GetAll() ([]any, error)
		CheckReflectionOptions
	}

	// Storage is the main cache interface that is used by the injector to store the values
	Storage[T comparable] interface {
		ValuesSetter[T]
		ValuesGetter[T]
	}

	// DependencyRetriever is the main element used to obtain registered binds/instances
	DependencyRetriever interface {
		ValuesGetter[BindKey]
		WrapRetriever() Injector
	}

	// Injector is the main interface that contains all needed methods to make an injector work
	Injector interface {
		Bind(BindKey, any) error
		BindNamed(BindKey, string, any) error
		SubInjector(allowOverrides ...bool) Injector
		DependencyRetriever
	}
)
