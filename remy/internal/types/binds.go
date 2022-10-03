package types

type (
	BindType      uint8
	Binder[T any] func(DependencyRetriever) T
	Bind[T any]   interface {
		Generates(DependencyRetriever) T
		Type() BindType
	}
)

const (
	BindInstance BindType = iota
	BindFactory
	BindSingleton
	BindLazySingleton
)
