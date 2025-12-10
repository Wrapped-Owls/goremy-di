package binds

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

type bindWrapper[T any] struct {
	types.Bind[T]
}

func (b bindWrapper[T]) PointerValue() any {
	return (*T)(nil)
}

func (b bindWrapper[T]) DefaultValue() any {
	var defaultValue T
	return defaultValue
}

func (b bindWrapper[T]) GenAsAny(retriever types.DependencyRetriever) (any, error) {
	return b.Generates(retriever)
}

func Singleton[T any](binder types.Binder[T]) types.Bind[T] {
	return bindWrapper[T]{&SingletonBind[T]{binder: binder}}
}

func LazySingleton[T any](binder types.Binder[T]) types.Bind[T] {
	return bindWrapper[T]{&SingletonBind[T]{binder: binder, IsLazy: true}}
}

func Factory[T any](binder types.Binder[T]) types.Bind[T] {
	return bindWrapper[T]{FactoryBind[T]{binder: binder, IsFactory: true}}
}

func Instance[T any](element T) types.Bind[T] {
	return bindWrapper[T]{
		FactoryBind[T]{
			binder: func(retriever types.DependencyRetriever) (T, error) {
				return element, nil
			},
		},
	}
}
