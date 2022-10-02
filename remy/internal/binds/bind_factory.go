package binds

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

type FactoryBind[T any] struct {
	IsFactory bool
	binder    types.Binder[T]
}

func (b FactoryBind[T]) Generates(injector types.DependencyRetriever) T {
	return b.binder(injector)
}

func (b FactoryBind[T]) Type() types.BindType {
	if b.IsFactory {
		return types.BindFactory
	}
	return types.BindInstance
}
