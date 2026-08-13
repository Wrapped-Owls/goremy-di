package binds

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

type FactoryBind[T any] struct {
	elemType[T]
	binder    types.Binder[T]
	IsFactory bool
}

func (b FactoryBind[T]) Generates(injector types.DependencyRetriever) (T, error) {
	return b.binder(injector)
}

func (b FactoryBind[T]) Type() types.BindType {
	if b.IsFactory {
		return types.BindFactory
	}
	return types.BindInstance
}

func (b FactoryBind[T]) GenAsAny(injector types.DependencyRetriever) (any, error) {
	return b.Generates(injector)
}
