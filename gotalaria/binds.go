package gotalaria

import (
	"gotalaria/internal/binds"
	"gotalaria/internal/types"
)

func Singleton[T any](binder types.Binder[T]) types.Bind[T] {
	return binds.Singleton(binder)
}

func LazySingleton[T any](binder types.Binder[T]) types.Bind[T] {
	return binds.LazySingleton(binder)
}

func Factory[T any](binder types.Binder[T]) types.Bind[T] {
	return binds.Factory(binder)
}

func Instance[T any](binder types.Binder[T]) types.Bind[T] {
	return binds.Instance(binder)
}
