package gotalaria

import (
	"gotalaria/internal/injector"
	"gotalaria/internal/types"
)

func NewInjector() types.Injector {
	return injector.New()
}

func Register[T any](i types.Injector, bind types.Bind[T], keys ...string) {
	injector.Register[T](i, bind, keys...)
}

func RegisterInstance[T any](i types.Injector, value T, keys ...string) {
	injector.Register[T](
		i,
		Instance[T](func(types.DependencyRetriever) T {
			return value
		}),
		keys...,
	)
}

func Get[T any](i types.DependencyRetriever, keys ...string) T {
	return injector.Get[T](i, keys...)
}
