package types

type (
	BindType      uint8
	Binder[T any] func(DependencyRetriever) (T, error)
	Bind[T any]   interface {
		Generates(DependencyRetriever) (T, error)
		Type() BindType
	}
)

const (
	BindInstance BindType = iota
	BindFactory
	BindSingleton
	BindLazySingleton
)

type BindOptions struct {
	Tag          string
	SoftOverride bool
}
