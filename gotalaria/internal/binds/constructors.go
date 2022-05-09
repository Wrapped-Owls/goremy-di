package binds

import "github.com/wrapped-owls/fitpiece/gotalaria/internal/types"

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
	return InstanceBind[T]{
		binder:    binder,
		IsFactory: true,
	}
}

func Instance[T any](binder types.Binder[T]) types.Bind[T] {
	return InstanceBind[T]{
		binder: binder,
	}
}
