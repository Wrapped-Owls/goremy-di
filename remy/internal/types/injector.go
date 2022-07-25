package types

type (
	CheckReflectionOptions interface {
		ReflectOpts() ReflectionOptions
	}
	ValuesSetter[T comparable] interface {
		// Set uses the T key given to save the value.
		Set(T, any) error

		// SetNamed uses the T key and cacheKey given to store the value.
		SetNamed(T, string, any) error
		CheckReflectionOptions
	}
	ValuesGetter[T comparable] interface {
		// GetNamed search for a named element that was cached using the T value given and a string key
		GetNamed(T, string) (any, error)

		// Get search for a named element that was cached using the T value given
		Get(T) (any, error)
		CheckReflectionOptions
	}

	// Storage is the main cache interface that is used by the injector to store the values
	Storage[T comparable] interface {
		ValuesSetter[T]
		ValuesGetter[T]
	}

	// DependencyRetriever is the main element used to obtain registered binds/instances
	DependencyRetriever interface {
		// RetrieveBind searches for a cached bind with the given BindKey, and if found it, returns.
		// This method only looks for existing binds, so instances will not be searched by it.
		RetrieveBind(BindKey) (any, error)

		// RetrieveNamedBind searches for a cached bind with the given BindKey and cacheKey, if found it, returns.
		// This method only looks for existing binds, so instances will not be searched by it.
		RetrieveNamedBind(string, BindKey) (any, error)
		ValuesGetter[BindKey]
	}

	// Injector is the main interface that contains all needed methods to make an injector work
	Injector interface {
		Bind(BindKey, any) error
		BindNamed(string, BindKey, any) error
		ValuesSetter[BindKey]
		DependencyRetriever
	}
)
