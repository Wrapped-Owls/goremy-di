package gotalaria

import (
	"github.com/wrapped-owls/fitpiece/gotalaria/internal/binds"
	"github.com/wrapped-owls/fitpiece/gotalaria/internal/types"
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
