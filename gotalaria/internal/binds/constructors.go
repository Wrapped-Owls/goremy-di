package binds

import "gotalaria/internal/types"

func Singleton[T any](binder types.Binder[T]) types.Bind[T] {
	return InstanceBind[T]{
		binder: binder,
	}
}

func LazySingleton[T any](binder types.Binder[T]) types.Bind[T] {
	return InstanceBind[T]{
		binder: binder,
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
