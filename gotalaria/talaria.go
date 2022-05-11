package gotalaria

import (
	"github.com/wrapped-owls/talaria-di/gotalaria/internal/injector"
	"github.com/wrapped-owls/talaria-di/gotalaria/internal/types"
)

type (
	DependencyRetriever = types.DependencyRetriever
	Injector            = types.Injector

	// Bind is directly copy from types.Bind
	Bind[T any] interface {
		Generates(DependencyRetriever) T
	}
)

func NewInjector() Injector {
	return injector.New(false, false)
}

// Register must be called first, because the library doesn't support registering dependencies while get at same time.
// This is not supported in multithreading applications because it does not have race protection
func Register[T any](i Injector, bind types.Bind[T], keys ...string) {
	if i == nil {
		i = globalInjector()
	}
	injector.Register[T](i, bind, keys...)
}

// RegisterInstance directly generates an instance bind without needing to write it.
//
// Receives: Injector (required); value (required); key (optional)
func RegisterInstance[T any](i Injector, value T, keys ...string) {
	Register[T](
		i,
		Instance[T](func(DependencyRetriever) T {
			return value
		}),
		keys...,
	)
}

// RegisterSingleton directly generates a singleton bind without needing to write it.
//
// Receives: Injector (required); value (required); key (optional)
func RegisterSingleton[T any](i Injector, value T, keys ...string) {
	Register[T](
		i,
		Singleton[T](func(DependencyRetriever) T {
			return value
		}),
		keys...,
	)
}

func Get[T any](i DependencyRetriever, keys ...string) T {
	if i == nil {
		i = globalInjector()
	}
	return injector.Get[T](i, keys...)
}
