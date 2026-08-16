package types

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
		// GetAll returns every element stored under keyTag, whatever its type, or the
		// untagged ones when keyTag is empty
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
)
