package types

type (
	BindKey = string

	InstancePair[T any] struct {
		Value T
		Key   string
	}

	ReflectionOptions struct {
		GenerifyInterface bool
		UseReflectionType bool
	}
)
