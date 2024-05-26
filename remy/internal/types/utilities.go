package types

type (
	InstancePair[T any] struct {
		// Value that will be injected directly. (required)
		Value T

		// Key must be used when registering the value as a bind
		Key string

		// InterfaceValue is just a pointer to an interface that will be used to register the value properly.
		//
		// It must only be used as the follow: `(*fmt.Stringer)(nil)`.
		//
		// As the final result, the bind will be registered
		// without the pointer. If you want to register an interface pointer, this option should not be used,
		// and you should pass the pointer to the Value directly.
		InterfaceValue any
	}

	ReflectionOptions struct {
		GenerifyInterface bool
		UseReflectionType bool
	}
)
