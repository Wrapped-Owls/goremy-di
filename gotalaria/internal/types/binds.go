package types

type (
	Binder[T any] func(DependencyRetriever) T
	Bind[T any]   interface {
		Generates(DependencyRetriever) (T, string)
		// TODO: Add a method Key() string
	}
)
