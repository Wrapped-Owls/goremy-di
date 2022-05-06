package gotalaria

import (
	"gotalaria/internal/injector"
	"gotalaria/internal/types"
)

func NewInjector() types.Injector {
	return injector.New()
}

func Register[T any](ij types.Injector, bind types.Bind[T], keys ...string) {
	injector.Register[T](ij, bind, keys...)
}

func Get[T any](ij types.DependencyRetriever, keys ...string) T {
	return injector.Get[T](ij, keys...)
}
