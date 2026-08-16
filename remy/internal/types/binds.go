package types

type (
	BindType      uint8
	Binder[T any] func(DependencyRetriever) (T, error)

	// Decorator enriches the value a bind produced before it reaches the caller
	Decorator[T any] func(DependencyRetriever, T) (T, error)
	Bind[T any]      interface {
		Generates(DependencyRetriever) (T, error)
		Type() BindType
	}

	// AnyGenerator generates an element without knowing its type parameter
	AnyGenerator interface {
		GenAsAny(retriever DependencyRetriever) (any, error)
	}

	// AnyBind is the type-erased view of a stored bind
	AnyBind interface {
		ElementKey() BindKey
		AnyGenerator
	}

	// GuessableBind tells whether a bind satisfies an interface before generating it
	GuessableBind interface {
		PointerValue() any
		DefaultValue() any
		AnyGenerator
	}
)

const (
	BindInstance BindType = iota
	BindFactory
	BindSingleton
	BindLazySingleton
)
