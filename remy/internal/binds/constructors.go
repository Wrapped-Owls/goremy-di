package binds

import "github.com/wrapped-owls/goremy-di/remy/internal/types"

func Singleton[T any](binder types.Binder[T]) types.Bind[T] {
	return &SingletonBind[T]{
		binder: binder,
	}
}

func LazySingleton[T any](binder types.Binder[T]) types.Bind[T] {
	return &SingletonBind[T]{
		binder: binder,
		IsLazy: true,
	}
}

func Factory[T any](binder types.Binder[T]) types.Bind[T] {
	return FactoryBind[T]{
		binder:    binder,
		IsFactory: true,
	}
}

func Instance[T any](binder types.Binder[T]) types.Bind[T] {
	return FactoryBind[T]{
		binder: binder,
	}
}
