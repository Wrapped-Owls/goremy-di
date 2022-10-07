package types

type (
	BindKey string

	InstancePair[T any] struct {
		Value       T
		Key         string
		IsInterface bool
	}

	ReflectionOptions struct {
		GenerifyInterface bool
		UseReflectionType bool
	}

	BindDependencies[T any] map[BindKey]T
	DependencyGraph         struct {
		UnnamedDependency BindDependencies[bool]
		NamedDependency   BindDependencies[map[string]bool]
	}
)
